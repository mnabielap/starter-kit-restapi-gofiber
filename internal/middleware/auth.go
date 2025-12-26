package middleware

import (
	"strings"

	"starter-kit-restapi-gofiber/internal/config"
	"starter-kit-restapi-gofiber/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Auth(cfg *config.Config, requiredRights ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please authenticate"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := utils.ValidateToken(tokenStr, cfg.JWTSecret)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please authenticate"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["type"] != utils.TokenTypeAccess {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token type"})
		}

		userID, _ := uuid.Parse(claims["sub"].(string))
		c.Locals("userId", userID)

		return c.Next()
	}
}

// RoleGuard checks if user has specific role (used after Auth)
func RoleGuard(getRole func(uuid.UUID) string, allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid := c.Locals("userId").(uuid.UUID)
		userRole := getRole(uid)
		
		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Forbidden"})
		}
		return c.Next()
	}
}