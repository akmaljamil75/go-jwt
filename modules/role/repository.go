package role

import (
	"go-jwt/common/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Save(role Role) (Role, error)
	FindOneRoleByName(name string) (Role, error)
	FindOneRoleByID(id uint) (Role, error)
	UpdateOne(role Role, input UpdateInputRole) (Role, error)
	FindOneAndLockAndUpdate(id uint, input UpdateInputRole) (Role, error)
	FindRolesByCrtieria(role Role) ([]Role, error)
	SoftDelete(id uint, input SoftDeleteInputRole) error
	RestoreSoftDelete(name string) (Role, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Save(role Role) (Role, error) {
	result := r.db.Save(&role)
	if result.Error != nil {
		return Role{}, result.Error
	}
	return role, nil
}

func (r *repository) FindOneRoleByName(name string) (Role, error) {
	var role = Role{Name: name}
	if err := r.db.Where(&role).First(&role); err.Error != nil {
		return Role{}, err.Error
	}
	return role, nil
}

func (r *repository) FindOneRoleByID(id uint) (Role, error) {
	var role = Role{ID: id}
	if err := r.db.Where(&role).First(&role).Error; err != nil {
		return Role{}, err
	}
	return role, nil
}

func (r *repository) UpdateOne(role Role, input UpdateInputRole) (Role, error) {
	if err := r.db.Model(&role).Updates(input).Error; err != nil {
		return Role{}, err
	}
	return role, nil
}

func (r *repository) FindOneAndLockAndUpdate(id uint, input UpdateInputRole) (Role, error) {

	var role Role

	err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).First(&role, id).Error; err != nil {
			return err
		}

		if role.Version != input.Version {
			return &response.FailedResponseMessage{
				Message: "Version mismatch",
				Code:    fiber.StatusConflict,
				Status:  "failed",
				Errors:  "The version of the resource you're trying to update has changed. Please make sure to get the latest version before trying again.",
			}
		}

		input.Version = time.Now().UnixMilli()

		if err := tx.Model(&role).Updates(input).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return Role{}, err
	}

	return role, nil
}

func (r *repository) FindRolesByCrtieria(role Role) ([]Role, error) {
	var result []Role
	if err := r.db.Where(&role).Find(&result).Error; err != nil {
		return []Role{}, err
	}
	return result, nil
}

func (r *repository) SoftDelete(id uint, input SoftDeleteInputRole) error {

	err := r.db.Transaction(func(tx *gorm.DB) error {

		var role Role

		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).First(&role, id).Error; err != nil {
			return err
		}

		if role.Version != input.Version {
			return &response.FailedResponseMessage{
				Message: "Version mismatch",
				Code:    fiber.StatusConflict,
				Status:  "failed",
				Errors:  "The version of the resource you're trying to update has changed. Please make sure to get the latest version before trying again.",
			}
		}

		if result := tx.Delete(&Role{ID: id}); result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) RestoreSoftDelete(name string) (Role, error) {

	err := r.db.Unscoped().Model(&Role{}).Where("name = ?", name).Update("deleted_at", nil)
	if err.Error != nil {
		return Role{}, err.Error
	}

	if err.RowsAffected == 0 {
		return Role{}, gorm.ErrRecordNotFound
	}

	var role Role
	if err := r.db.Where(&Role{Name: name}).First(&role).Error; err != nil {
		return Role{}, err
	}

	return role, nil
}
