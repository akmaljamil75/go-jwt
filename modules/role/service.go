package role

import (
	"errors"
	"go-jwt/common/response"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type Service interface {
	Save(input RegisterInputRole) (Role, response.FailedResponseMessage)
	FindOneRoleByName(name string) (Role, response.FailedResponseMessage)
	UpdateOne(id uint, input UpdateInputRole) (Role, response.FailedResponseMessage)
	FindOneRoleByID(id uint) (Role, response.FailedResponseMessage)
	FindRolesByCrtieria(role Role) ([]Role, response.FailedResponseMessage)
	SoftDelete(id uint, input SoftDeleteInputRole) response.FailedResponseMessage
	RestoreDataSoftDelete(name string) (Role, response.FailedResponseMessage)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Save(input RegisterInputRole) (Role, response.FailedResponseMessage) {

	save, err := s.repo.Save(Role{Name: input.Name, Version: time.Now().UnixMilli()})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return Role{}, response.FailedResponseMessage{
				Message: "Duplicated key for role " + input.Name,
				Status:  "failed",
				Code:    http.StatusBadRequest,
				Errors:  err.Error(),
			}
		} else {
			return Role{}, response.FailedResponseMessage{
				Message: "Failed to save role",
				Status:  "failed",
				Code:    http.StatusInternalServerError,
				Errors:  err.Error(),
			}
		}

	}
	return save, response.FailedResponseMessage{}
}

func (s *service) FindOneRoleByName(name string) (Role, response.FailedResponseMessage) {
	role, err := s.repo.FindOneRoleByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return Role{}, response.FailedResponseMessage{
				Message: "Role not found",
				Status:  "failed",
				Code:    http.StatusNotFound,
				Errors:  err.Error(),
			}
		} else {
			return Role{}, response.FailedResponseMessage{
				Message: "Failed to get role",
				Status:  "failed",
				Code:    http.StatusInternalServerError,
				Errors:  err.Error(),
			}
		}
	}
	return role, response.FailedResponseMessage{}
}

func (s *service) UpdateOne(id uint, input UpdateInputRole) (Role, response.FailedResponseMessage) {
	role, err := s.repo.FindOneAndLockAndUpdate(id, input)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return Role{}, response.FailedResponseMessage{
				Message: "Duplicated key for role " + input.Name,
				Status:  "failed",
				Code:    http.StatusBadRequest,
				Errors:  err.Error(),
			}
		}
		return Role{}, response.FailedResponseMessage{
			Message: "Failed to update role",
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Errors:  err.Error(),
		}
	}
	return role, response.FailedResponseMessage{}
}

func (s *service) FindOneRoleByID(id uint) (Role, response.FailedResponseMessage) {
	role, err := s.repo.FindOneRoleByID(id)
	if err != nil {

		var responseErr *response.FailedResponseMessage
		if errors.As(err, &responseErr) {
			return Role{}, *responseErr
		} else if err == gorm.ErrRecordNotFound {
			return Role{}, response.FailedResponseMessage{
				Message: "Role not found",
				Status:  "failed",
				Code:    http.StatusNotFound,
				Errors:  err.Error(),
			}
		} else {
			return Role{}, response.FailedResponseMessage{
				Message: "Failed to get role",
				Status:  "failed",
				Code:    http.StatusInternalServerError,
				Errors:  err.Error(),
			}
		}
	}
	return role, response.FailedResponseMessage{}
}

func (s *service) FindRolesByCrtieria(role Role) ([]Role, response.FailedResponseMessage) {

	var roles []Role

	roles, err := s.repo.FindRolesByCrtieria(role)

	if err != nil {
		return []Role{}, response.FailedResponseMessage{
			Message: "failed to find role",
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Errors:  err.Error(),
		}
	}

	return roles, response.FailedResponseMessage{}
}

// SoftDelete implements Service.
func (s *service) SoftDelete(id uint, input SoftDeleteInputRole) response.FailedResponseMessage {

	if err := s.repo.SoftDelete(id, input); err != nil {
		var responseErr *response.FailedResponseMessage
		if errors.As(err, &responseErr) {
			return *responseErr
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.FailedResponseMessage{
				Message: "Role not found",
				Status:  "failed",
				Code:    http.StatusNotFound,
				Errors:  err.Error(),
			}
		} else {
			return response.FailedResponseMessage{
				Message: "Failed to soft delete role",
				Status:  "failed",
				Code:    http.StatusInternalServerError,
				Errors:  err.Error(),
			}
		}
	}

	return response.FailedResponseMessage{}
}

func (s *service) RestoreDataSoftDelete(name string) (Role, response.FailedResponseMessage) {

	role, err := s.repo.RestoreSoftDelete(name)
	if err != nil {

		var responseErr *response.FailedResponseMessage
		if errors.As(err, &responseErr) {
			return Role{}, *responseErr
		} else if err == gorm.ErrRecordNotFound {
			return Role{}, response.FailedResponseMessage{
				Message: "Role not found",
				Status:  "failed",
				Code:    http.StatusNotFound,
				Errors:  err.Error(),
			}
		} else {
			return Role{}, response.FailedResponseMessage{
				Message: "Failed to restore role",
				Status:  "failed",
				Code:    http.StatusInternalServerError,
				Errors:  err.Error(),
			}
		}

	}

	return role, response.FailedResponseMessage{}
}
