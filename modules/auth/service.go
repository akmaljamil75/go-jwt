package auth

import (
	"errors"
	"go-jwt/common/jwt"
	"go-jwt/common/response"
	"go-jwt/modules/role"
	"go-jwt/modules/user"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Login(username, password string) (string, response.FailedResponseMessage)
	VertifikasiToken(token string) response.FailedResponseMessage
}

type service struct {
	userRepo user.Repository
	roleRepo role.Repository
}

// VertifikasiToken implements Service.

func NewService(uRepo user.Repository, rRepo role.Repository) Service {
	return &service{uRepo, rRepo}
}

func (s *service) Login(username string, password string) (string, response.FailedResponseMessage) {

	user, err := s.userRepo.FindUserOneUserByUsername(username)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", response.FailedResponseMessage{
				Message: "Invalid username or password",
				Status:  "failed",
				Code:    http.StatusUnauthorized,
				Errors:  err.Error(),
			}
		}

		return "", response.FailedResponseMessage{
			Message: "failed to find username",
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Errors:  err.Error(),
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {

		if err == bcrypt.ErrMismatchedHashAndPassword {
			return "", response.FailedResponseMessage{
				Message: "Invalid username or password",
				Status:  "failed",
				Code:    http.StatusUnauthorized,
				Errors:  "Invalid username or password",
			}
		}

		return "", response.FailedResponseMessage{
			Message: "Role not found",
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Errors:  err.Error(),
		}
	}

	role, err := s.roleRepo.FindOneRoleByID(user.RoleID)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", response.FailedResponseMessage{
				Message: "Role not found",
				Status:  "failed",
				Code:    http.StatusBadRequest,
				Errors:  err.Error(),
			}
		}

		return "", response.FailedResponseMessage{
			Message: "Failed to find Role",
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Errors:  err.Error(),
		}
	}

	token, err := jwt.GenerateToken(user.Username, role.ID)
	if err != nil {
		return "", response.FailedResponseMessage{
			Message: "Failed to generate token",
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Errors:  err.Error(),
		}
	}
	return token, response.FailedResponseMessage{}
}

func (s *service) VertifikasiToken(token string) response.FailedResponseMessage {

	_, err := jwt.VerifyToken(token)
	if err != nil {

		var responseMessageFailed *response.FailedResponseMessage
		if errors.As(err, &responseMessageFailed) {
			return *responseMessageFailed
		}

		return response.FailedResponseMessage{
			Message: "Failed to Verify token",
			Status:  "failed",
			Code:    fiber.StatusInternalServerError,
			Errors:  err.Error(),
		}
	}

	return response.FailedResponseMessage{}
}
