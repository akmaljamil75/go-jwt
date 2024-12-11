package middleware

import (
	"context"
	"errors"
	"go-jwt/common/response"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func HandlingErrorMiddleware(c *fiber.Ctx) error {

	err := c.Next()

	var responseTemplate *response.FailedResponseMessage

	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return c.Status(fiber.StatusGatewayTimeout).JSON(response.FailedResponseMessage{
				Message: "Request timed out",
				Status:  "failed",
				Code:    fiber.StatusGatewayTimeout,
			})
		case errors.As(err, &responseTemplate):
			return c.Status(responseTemplate.Code).JSON(responseTemplate)

		case errors.Is(err, gorm.ErrRecordNotFound):
			return c.Status(fiber.StatusInternalServerError).JSON(response.FailedResponseMessage{
				Message: "Record not found",
				Status:  "failed",
				Code:    fiber.StatusInternalServerError,
			})
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return c.Status(fiber.StatusInternalServerError).JSON(response.FailedResponseMessage{
				Message: "Duplicated key " + err.Error(),
				Status:  "failed",
				Code:    fiber.StatusInternalServerError,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(response.FailedResponseMessage{
				Message: "Internal Server Error",
				Status:  "failed",
				Code:    fiber.StatusInternalServerError,
			})
		}
	}
	return nil
}
