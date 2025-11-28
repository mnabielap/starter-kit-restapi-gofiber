package routes

import (
	"starter-kit-restapi-gofiber/internal/config"
	"starter-kit-restapi-gofiber/internal/handlers"
	"starter-kit-restapi-gofiber/internal/services"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	// 1. Initialize Services
	tokenService := services.NewTokenService(db, cfg)
	userService := services.NewUserService(db)
	emailService := services.NewEmailService(cfg)
	authService := services.NewAuthService(userService, tokenService, emailService)

	// 2. Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// 3. Swagger Route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// 4. Group v1 API
	v1 := app.Group("/v1")

	// 5. Register Module Routes
	SetupAuthRoutes(v1, authHandler, cfg)
	SetupUserRoutes(v1, userHandler, cfg, userService)
}