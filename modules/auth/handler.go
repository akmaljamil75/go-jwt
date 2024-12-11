package auth

import (
	"go-jwt/common/response"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) Login(c *fiber.Ctx) error {

	var input LoginInput

	if err := c.BodyParser(&input); err != nil {
		return &response.FailedResponseMessage{
			Message: "Failed to parse request body",
			Status:  "failed",
			Code:    fiber.StatusUnprocessableEntity,
			Errors:  err.Error(),
		}
	}

	validation := response.ValidateBodyRequest(input)
	if len(validation) != 0 {
		return &response.FailedResponseMessage{
			Message: "Failed request body",
			Status:  "failed",
			Code:    fiber.StatusBadRequest,
			Errors:  validation,
		}
	}

	token, err := h.service.Login(input.Username, input.Password)
	if !reflect.DeepEqual(err, response.FailedResponseMessage{}) {
		return &err
	}

	return c.Status(fiber.StatusOK).JSON(response.BuildSuccessResponseMessage("login successfully", 200, map[string]string{
		"access_token": token,
	}))
}
