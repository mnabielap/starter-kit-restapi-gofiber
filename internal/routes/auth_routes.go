package routes

import (
	"time"

	"starter-kit-restapi-gofiber/internal/config"
	"starter-kit-restapi-gofiber/internal/handlers"
	"starter-kit-restapi-gofiber/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func SetupAuthRoutes(router fiber.Router, h *handlers.AuthHandler, cfg *config.Config) {
	auth := router.Group("/auth")

	// Rate Limiter specifically for Auth endpoints
	authLimiter := limiter.New(limiter.Config{
		Max:        20,
		Expiration: 15 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many requests, please try again later.",
			})
		},
	})
	auth.Use(authLimiter)

	// Public Routes
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/logout", h.Logout)
	auth.Post("/refresh-tokens", h.RefreshTokens)
	auth.Post("/forgot-password", h.ForgotPassword)
	auth.Post("/reset-password", h.ResetPassword)
	auth.Post("/verify-email", h.VerifyEmail)

	// Protected Routes (Requires JWT)
	auth.Post("/send-verification-email", middleware.Auth(cfg), h.SendVerificationEmail)
}