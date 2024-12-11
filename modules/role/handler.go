package role

import (
	"go-jwt/common/response"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service}
}

func (h *handler) Create(c *fiber.Ctx) error {

	var input RegisterInputRole

	if err := c.BodyParser(&input); err != nil {
		return &response.FailedResponseMessage{
			Message: "Failed to parse request body",
			Status:  "failed",
			Code:    fiber.StatusUnprocessableEntity,
			Errors:  err.Error(),
		}
	}
	user, err := h.service.Save(input)
	if !reflect.DeepEqual(err, response.FailedResponseMessage{}) {
		return &err
	}

	return c.Status(fiber.StatusOK).JSON(response.BuildSuccessResponseMessage("successfully created role", http.StatusOK, user))
}

func (h *handler) FindOneRoleByName(c *fiber.Ctx) error {
	name := c.Params("name")
	user, err := h.service.FindOneRoleByName(name)
	if !reflect.DeepEqual(err, response.FailedResponseMessage{}) {
		return &err
	}
	return c.Status(fiber.StatusOK).JSON(response.BuildSuccessResponseMessage("successfully find role", http.StatusOK, user))
}

func (h *handler) FindOneRoleByID(c *fiber.Ctx) error {

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return &response.FailedResponseMessage{
			Message: "Invalid Convert ID",
			Status:  "failed",
			Code:    fiber.StatusInternalServerError,
			Errors:  err.Error(),
		}
	}

	uintID := uint(id)
	user, errFind := h.service.FindOneRoleByID(uintID)
	if !reflect.DeepEqual(err, response.FailedResponseMessage{}) {
		return &errFind
	}
	return c.Status(fiber.StatusOK).JSON(response.BuildSuccessResponseMessage("successfully find role", http.StatusOK, user))
}

func (h *handler) Update(c *fiber.Ctx) error {

	var input UpdateInputRole
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

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return &response.FailedResponseMessage{
			Message: "Invalid Convert ID",
			Status:  "failed",
			Code:    fiber.StatusInternalServerError,
			Errors:  err.Error(),
		}
	}
	uintID := uint(id)

	update, errUpdate := h.service.UpdateOne(uintID, input)
	if !reflect.DeepEqual(errUpdate, response.FailedResponseMessage{}) {
		return &errUpdate
	}

	return c.Status(fiber.StatusOK).JSON(response.BuildSuccessResponseMessage("successfully updated role", http.StatusOK, update))
}

func (h *handler) FindRoles(c *fiber.Ctx) error {

	var criteria Role
	var roles []Role

	if err := c.BodyParser(&criteria); err != nil {
		return &response.FailedResponseMessage{
			Message: "Failed to parse request body",
			Status:  "failed",
			Code:    fiber.StatusUnprocessableEntity,
			Errors:  err.Error(),
		}
	}

	roles, err := h.service.FindRolesByCrtieria(criteria)
	if !reflect.DeepEqual(err, response.FailedResponseMessage{}) {
		return &err
	}

	return c.Status(fiber.StatusOK).JSON(response.BuildSuccessResponseMessage("successfully find role", http.StatusOK, roles))
}

func (h *handler) SoftDelete(c *fiber.Ctx) error {

	var input SoftDeleteInputRole
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

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return &response.FailedResponseMessage{
			Message: "Invalid Convert ID",
			Status:  "failed",
			Code:    fiber.StatusInternalServerError,
			Errors:  err.Error(),
		}
	}
	uintID := uint(id)

	errUpdate := h.service.SoftDelete(uintID, input)
	if !reflect.DeepEqual(errUpdate, response.FailedResponseMessage{}) {
		return &errUpdate
	}

	return c.Status(fiber.StatusOK).JSON(response.BuildSuccessResponseMessage("successfully soft deleted role", http.StatusOK, nil))
}

func (h *handler) RestoreSoftDelete(c *fiber.Ctx) error {

	name := c.Params("name")

	role, err := h.service.RestoreDataSoftDelete(name)
	if !reflect.DeepEqual(err, response.FailedResponseMessage{}) {
		return &err
	}

	return c.Status(fiber.StatusOK).JSON(response.BuildSuccessResponseMessage("successfully restored role", http.StatusOK, role))
}
