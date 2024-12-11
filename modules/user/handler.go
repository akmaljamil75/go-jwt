package user

import (
	"go-jwt/common/response"
	"net/http"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type handler struct {
	userService Service
}

func NewHandler(userService Service) *handler {
	return &handler{userService}
}

func (h *handler) Create(c *fiber.Ctx) error {

	var input RegisterInputUser

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

	user, err := h.userService.Save(input)
	if !reflect.DeepEqual(err, response.FailedResponseMessage{}) {
		return &err
	}

	return c.Status(fiber.StatusOK).JSON(response.BuildSuccessResponseMessage("successfuly created user", http.StatusOK, user))
}

func (h *handler) SoftDelete(c *fiber.Ctx) error {

	var input SoftDeleteInputUser

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

	if err := h.userService.SoftDelete(input.ID, input.Version); !reflect.DeepEqual(err, response.FailedResponseMessage{}) {
		return &err
	}

	return c.Status(fiber.StatusOK).JSON(response.BuildSuccessResponseMessage("successfully deleted user", http.StatusOK, nil))
}

func (h *handler) FindOneByUsername(c *fiber.Ctx) error {
	username := c.Params("username")

	user, err := h.userService.FindOneUserByUsername(username)
	if !reflect.DeepEqual(err, response.FailedResponseMessage{}) {
		return &err
	}

	return c.Status(fiber.StatusOK).JSON(response.BuildSuccessResponseMessage("successfully find user", http.StatusOK, user))
}
