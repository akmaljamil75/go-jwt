package jwt

import (
	"errors"
	"go-jwt/common/database"
	"go-jwt/common/response"
	"go-jwt/modules/role"
	"go-jwt/modules/user"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var JWT_SIGNATURE_KEY = []byte("the secret of kalimdor")

func GenerateToken(username string, roleID uint) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role_id":  roleID,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
		"issuer":   "go-jwt",
		"aud":      "go-jwt-client",
	}, func(t *jwt.Token) {
		t.Header["typ"] = "JWT"
	})

	tokenString, err := token.SignedString(JWT_SIGNATURE_KEY)
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &response.FailedResponseMessage{
				Message: "Invalid signing method",
				Status:  "failed",
				Code:    fiber.StatusUnauthorized,
				Errors:  nil,
			}
		}
		return JWT_SIGNATURE_KEY, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, &response.FailedResponseMessage{
				Message: "token expired",
				Status:  "failed",
				Code:    fiber.StatusUnauthorized,
				Errors:  nil,
			}
		}
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, &response.FailedResponseMessage{
			Message: "invalid token",
			Status:  "failed",
			Code:    fiber.StatusUnauthorized,
			Errors:  nil,
		}
	}

	username := (claims)["username"].(string)
	roleID := (claims)["role_id"].(float64)
	exp := (claims)["exp"].(float64)
	issuer := (claims)["issuer"].(string)
	aud := (claims)["aud"].(string)

	if time.Now().Unix() > int64(exp) {
		return nil, &response.FailedResponseMessage{
			Message: "token expired",
			Status:  "failed",
			Code:    fiber.StatusUnauthorized,
			Errors:  nil,
		}
	} else if issuer != "go-jwt" || aud != "go-jwt-client" {
		return nil, &response.FailedResponseMessage{
			Message: "invalid token",
			Status:  "failed",
			Code:    fiber.StatusUnauthorized,
			Errors:  nil,
		}
	} else {

		db := database.GetDB()

		var user user.User
		if err := db.First(&user, "username = ?", username).Error; err != nil {
			return nil, &response.FailedResponseMessage{
				Message: "invalid username",
				Status:  "failed",
				Code:    fiber.StatusUnauthorized,
				Errors:  err.Error(),
			}
		}
		var role role.Role
		if err := db.First(&role, roleID).Error; err != nil {
			return nil, &response.FailedResponseMessage{
				Message: "invalid role id",
				Status:  "failed",
				Code:    fiber.StatusUnauthorized,
				Errors:  err.Error(),
			}
		}

	}

	return &claims, nil
}
