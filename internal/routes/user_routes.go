package routes

import (
	"starter-kit-restapi-gofiber/internal/config"
	"starter-kit-restapi-gofiber/internal/handlers"
	"starter-kit-restapi-gofiber/internal/middleware"
	"starter-kit-restapi-gofiber/internal/services"
	"starter-kit-restapi-gofiber/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SetupUserRoutes(router fiber.Router, h *handlers.UserHandler, cfg *config.Config, userService *services.UserService) {
	users := router.Group("/users")

	// 1. Apply Auth Middleware to ALL user routes
	users.Use(middleware.Auth(cfg))

	// 2. Define Admin Guard
	// This helper function fetches the user role from DB to check permissions
	adminGuard := middleware.RoleGuard(func(uid uuid.UUID) string {
		user, err := userService.GetUserById(uid.String())
		if err != nil {
			return ""
		}
		return user.Role
	}, utils.RoleAdmin)

	// 3. Admin Only Routes
	users.Post("/", adminGuard, h.CreateUser)
	users.Get("/", adminGuard, h.GetUsers)

	// 4. User/Admin Routes (Logic for self-access is handled inside Handler)
	users.Get("/:userId", h.GetUser)
	users.Patch("/:userId", h.UpdateUser)
	users.Delete("/:userId", h.DeleteUser)
}