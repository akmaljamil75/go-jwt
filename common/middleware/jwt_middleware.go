package middleware

import (
	"errors"
	"go-jwt/common/jwt"
	"go-jwt/common/response"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JwtAuthorization(c *fiber.Ctx) error {

	auth := c.Get("Authorization")

	if auth == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(response.FailedResponseMessage{Code: fiber.StatusUnauthorized, Message: "Missing Authorization Header", Status: "failed", Errors: "Missing Authorization Header"})
	}

	splitToken := strings.Split(auth, "Bearer ")
	var token string
	if len(splitToken) > 1 {
		token = splitToken[1]
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(response.FailedResponseMessage{Code: fiber.StatusUnauthorized, Message: "Invalid Authorization Header", Status: "failed", Errors: "Invalid Authorization Header"})
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(response.FailedResponseMessage{Code: fiber.StatusUnauthorized, Message: "Invalid token", Status: "failed", Errors: "Invalid token"})
	}

	claims, err := jwt.VerifyToken(token)
	if err != nil {

		var responseErr *response.FailedResponseMessage

		if errors.As(err, &responseErr) {
			return responseErr
		}

		return c.Status(fiber.StatusUnauthorized).JSON(response.FailedResponseMessage{Code: fiber.StatusUnauthorized, Message: "Invalid token", Status: "failed", Errors: err.Error()})
	}

	c.Locals("claims", claims)
	return c.Next()
}
