package user

import (
	"errors"
	"go-jwt/common/response"
	"go-jwt/modules/role"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Save(input RegisterInputUser) (User, response.FailedResponseMessage)
	FindOneUserByUsername(username string) (User, response.FailedResponseMessage)
	FindUsersByCriteria(user User) ([]User, response.FailedResponseMessage)
	SoftDelete(id uint, version int64) response.FailedResponseMessage
	Update(id uint, input UpdateInputUser) (User, response.FailedResponseMessage)
}

type service struct {
	userRepo Repository
	roleRepo role.Repository
}

func NewService(userRepo Repository, roleRepo role.Repository) Service {
	return &service{userRepo: userRepo, roleRepo: roleRepo}
}

func (s *service) Save(input RegisterInputUser) (User, response.FailedResponseMessage) {

	role, err := s.roleRepo.FindOneRoleByID(input.RoleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return User{}, response.FailedResponseMessage{
				Message: "role not found",
				Status:  "failed",
				Code:    http.StatusBadRequest,
				Errors:  err.Error(),
			}
		} else {
			return User{}, response.FailedResponseMessage{
				Message: "failed to find role by id for check role is empty or not empty",
				Status:  "failed",
				Code:    http.StatusInternalServerError,
				Errors:  err.Error(),
			}
		}
	}

	if role.ID != input.RoleID {
		return User{}, response.FailedResponseMessage{
			Message: "role not found",
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Errors:  "role not found",
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, response.FailedResponseMessage{
			Message: "Failed to hash password",
			Status:  "failed",
			Errors:  err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}

	var toSaveUser = User{
		Username: input.Username,
		RoleID:   role.ID,
		Password: string(passwordHash),
		Version:  time.Now().UnixMilli(),
	}

	user, err := s.userRepo.Save(toSaveUser)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return User{}, response.FailedResponseMessage{
				Message: "Duplicated key for username " + input.Username,
				Status:  "failed",
				Errors:  err.Error(),
				Code:    http.StatusBadRequest,
			}
		} else {
			return User{}, response.FailedResponseMessage{
				Message: "Failed to save user",
				Status:  "failed",
				Errors:  err.Error(),
				Code:    http.StatusInternalServerError,
			}
		}
	}
	return user, response.FailedResponseMessage{}
}

// FindUserOneUserByUsername implements Service.
func (s *service) FindOneUserByUsername(username string) (User, response.FailedResponseMessage) {

	user, err := s.userRepo.FindUserOneUserByUsername(username)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, response.FailedResponseMessage{
				Message: "User not found",
				Status:  "failed",
				Code:    http.StatusNotFound,
				Errors:  err.Error(),
			}
		} else {
			return User{}, response.FailedResponseMessage{
				Message: "Failed to find user by username",
				Status:  "failed",
				Code:    http.StatusInternalServerError,
				Errors:  err.Error(),
			}
		}

	}

	return user, response.FailedResponseMessage{}
}

// FindUsersByCriteria implements Service.
func (s *service) FindUsersByCriteria(user User) ([]User, response.FailedResponseMessage) {

	users, err := s.userRepo.FindUsersByCriteria(user)
	if err != nil {
		return users, response.FailedResponseMessage{
			Message: "Failed to find users by criteria",
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Errors:  err.Error(),
		}
	}

	return users, response.FailedResponseMessage{}
}

func (s *service) SoftDelete(id uint, version int64) response.FailedResponseMessage {
	if err := s.userRepo.SoftDelete(id, version); err != nil {
		var failedResponse *response.FailedResponseMessage
		if errors.As(err, &failedResponse) {
			return *failedResponse
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.FailedResponseMessage{
				Message: "Record not found",
				Status:  "failed",
				Code:    http.StatusNotFound,
				Errors:  err.Error(),
			}
		} else {
			return response.FailedResponseMessage{
				Message: "Failed to soft delete user",
				Status:  "failed",
				Code:    http.StatusInternalServerError,
				Errors:  err.Error(),
			}
		}
	}
	return response.FailedResponseMessage{}
}

func (s *service) Update(id uint, input UpdateInputUser) (User, response.FailedResponseMessage) {

	user, err := s.userRepo.UpdateOne(id, input)
	if err != nil {
		var responseFailed *response.FailedResponseMessage
		if errors.As(err, &responseFailed) {
			return User{}, *responseFailed
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, response.FailedResponseMessage{
				Message: "User not found",
				Status:  "failed",
				Code:    http.StatusNotFound,
				Errors:  err.Error(),
			}
		} else {
			return User{}, response.FailedResponseMessage{
				Message: "Failed to update user",
				Status:  "failed",
				Code:    http.StatusInternalServerError,
				Errors:  err.Error(),
			}
		}

	}
	return user, response.FailedResponseMessage{}
}
