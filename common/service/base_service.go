package service

import (
	"go-jwt/common/response"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service interface {
	Save(input interface{}) (interface{}, error)
	FindOne(criteria interface{}) (interface{}, error)
	Find(criteria interface{}) (interface{}, error)
	SoftDelete(id uint, version int64) error
	Update(id uint, input interface{}) (interface{}, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

func (s *service) Save(v interface{}) (interface{}, error) {
	result := s.db.Save(v)
	if result.Error != nil {
		return nil, result.Error
	}
	return v, nil
}

func (s *service) FindOne(criteria interface{}) (interface{}, error) {
	var result interface{}
	if err := s.db.Where(criteria).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) Find(criteria interface{}) (interface{}, error) {
	var result []interface{}
	if err := s.db.Where(criteria).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) SoftDelete(id uint, version int64) error {

	err := s.db.Transaction(func(tx *gorm.DB) error {

		var model interface{}
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).Where(&model, id).First(&model).Error; err != nil {
			return err
		}

		if reflect.ValueOf(model).Elem().FieldByName("Version").Int() != version {
			return &response.FailedResponseMessage{
				Message: "Version mismatch",
				Code:    fiber.StatusConflict,
				Status:  "failed",
			}
		}

		if err := tx.Delete(&model).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *service) Update(id uint, input interface{}) (interface{}, error) {

	var model interface{}

	err := s.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).Where(&model, id).First(&model).Error; err != nil {
			return err
		}

		if reflect.ValueOf(model).Elem().FieldByName("Version").Int() != reflect.ValueOf(input).Elem().FieldByName("Version").Int() {
			return &response.FailedResponseMessage{
				Message: "Version mismatch",
				Code:    fiber.StatusConflict,
				Status:  "failed",
			}
		}

		if err := tx.Model(&model).Updates(input).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return model, nil
}
