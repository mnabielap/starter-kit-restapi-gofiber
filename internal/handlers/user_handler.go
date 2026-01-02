package handlers

import (
	"starter-kit-restapi-gofiber/internal/dto"
	"starter-kit-restapi-gofiber/internal/services"
	"starter-kit-restapi-gofiber/pkg/utils"
	"starter-kit-restapi-gofiber/pkg/validator"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

// CreateUser
// @Tags Users
// @Security BearerAuth
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := validator.ParseAndValidate(c, &req); err != nil { return nil }
	user, err := h.Service.CreateUser(&req)
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()}) }
	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUsers
// @Tags Users
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	name := c.Query("name")
	role := c.Query("role")

	pagination := &utils.Pagination{Limit: limit, Page: page}
	result, err := h.Service.QueryUsers(pagination, name, role)
	if err != nil { return c.SendStatus(fiber.StatusInternalServerError) }
	return c.JSON(result)
}

// GetUser
// @Tags Users
// @Security BearerAuth
// @Router /users/{userId} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("userId")
	currentID := c.Locals("userId").(uuid.UUID)
	
	user, err := h.Service.GetUserById(id)
	if err != nil { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"}) }

	if user.ID != currentID {
		requester, err := h.Service.GetUserById(currentID.String())
		if err != nil || requester.Role != utils.RoleAdmin {
			return c.SendStatus(fiber.StatusForbidden)
		}
	}

	return c.JSON(user)
}

// UpdateUser
// @Tags Users
// @Security BearerAuth
// @Router /users/{userId} [patch]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("userId")
	currentID := c.Locals("userId").(uuid.UUID)

	// Authorization check: Prevent IDOR.
	// Users can only update themselves. Admins can update anyone.
	if id != currentID.String() {
		requester, err := h.Service.GetUserById(currentID.String())
		if err != nil || requester.Role != utils.RoleAdmin {
			return c.SendStatus(fiber.StatusForbidden)
		}
	}

	var req dto.UpdateUserRequest
	if err := validator.ParseAndValidate(c, &req); err != nil { return nil }
	
	user, err := h.Service.UpdateUserById(id, &req)
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()}) }
	return c.JSON(user)
}

// DeleteUser
// @Tags Users
// @Security BearerAuth
// @Router /users/{userId} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("userId")
	if err := h.Service.DeleteUserById(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}