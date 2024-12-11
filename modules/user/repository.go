package user

import (
	"go-jwt/common/response"
	"go-jwt/modules/role"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Save(user User) (User, error)
	FindUserOneUserByUsername(username string) (User, error)
	FindUsersByCriteria(user User) ([]User, error)
	SoftDelete(id uint, version int64) error
	UpdateOne(id uint, user UpdateInputUser) (User, error)
	FindOneRoleByUsername(username string) (role.Role, error)
}

type repository struct {
	db *gorm.DB
}

// Save implements Repository.
func (r *repository) Save(user User) (User, error) {
	result := r.db.Save(&user).Error
	if result != nil {
		return user, result
	}
	return user, nil
}

// FindOneRoleByUsername implements Repository.
func (r *repository) FindOneRoleByUsername(username string) (role.Role, error) {
	var role role.Role
	if err := r.db.Where(User{Username: username}).Preload(clause.Associations).First(&role).Error; err != nil {
		return role, err
	}
	return role, nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindUserOneUserByUsername(username string) (User, error) {
	var user User
	if err := r.db.Where(User{Username: username}).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *repository) FindUsersByCriteria(user User) ([]User, error) {
	var users []User

	if err := r.db.Where(&user).Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (r *repository) SoftDelete(id uint, version int64) error {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		var user User
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).First(&user, id).Error; err != nil {
			return err
		}

		if user.Version != version {
			return &response.FailedResponseMessage{
				Message: "Version mismatch",
				Code:    fiber.StatusConflict,
				Status:  "failed",
				Errors:  "The version of the resource you're trying to update has changed. Please make sure to get the latest version before trying again.",
			}
		}

		if err := tx.Delete(&user).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateOne(id uint, input UpdateInputUser) (User, error) {

	var user User

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).Where(&user, id).First(&user).Error; err != nil {
			return err
		}

		var role role.Role
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).First(&role, input.RoleID).Error; err != nil {
			return err
		}

		if user.Version != input.Version {
			return &response.FailedResponseMessage{
				Message: "Version mismatch",
				Code:    fiber.StatusConflict,
				Status:  "failed",
				Errors:  "The version of the resource you're trying to update has changed. Please make sure to get the latest version before trying again.",
			}
		}

		input.Version = time.Now().UnixMilli()

		result := r.db.Model(&user).Updates(input)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected > 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	}); err != nil {
		return user, nil
	}

	return user, nil
}
