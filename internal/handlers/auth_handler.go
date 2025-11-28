package handlers

import (
	"starter-kit-restapi-gofiber/internal/dto"
	"starter-kit-restapi-gofiber/internal/services"
	"starter-kit-restapi-gofiber/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	Service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{Service: s}
}

// Register
// @Summary Register
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.RegisterRequest true "Body"
// @Success 201 {object} dto.AuthResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := validator.ParseAndValidate(c, &req); err != nil { return nil }
	res, err := h.Service.Register(&req)
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()}) }
	return c.Status(fiber.StatusCreated).JSON(res)
}

// Login
// @Summary Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Body"
// @Success 200 {object} dto.AuthResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := validator.ParseAndValidate(c, &req); err != nil { return nil }
	res, err := h.Service.Login(&req)
	if err != nil { return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()}) }
	return c.Status(fiber.StatusOK).JSON(res)
}

// Logout
// @Summary Logout
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LogoutRequest true "Body"
// @Success 204
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req dto.LogoutRequest
	if err := validator.ParseAndValidate(c, &req); err != nil { return nil }
	h.Service.Logout(req.RefreshToken)
	return c.SendStatus(fiber.StatusNoContent)
}

// Refresh Tokens
// @Summary Refresh Tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.RefreshTokenRequest true "Body"
// @Success 200 {object} dto.AuthResponse
// @Router /auth/refresh-tokens [post]
func (h *AuthHandler) RefreshTokens(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := validator.ParseAndValidate(c, &req); err != nil { return nil }
	res, err := h.Service.RefreshAuth(req.RefreshToken)
	if err != nil { return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()}) }
	return c.JSON(res)
}

// Forgot Password
// @Summary Forgot Password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.ForgotPasswordRequest true "Body"
// @Success 204
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req dto.ForgotPasswordRequest
	if err := validator.ParseAndValidate(c, &req); err != nil { return nil }
	h.Service.ForgotPassword(req.Email)
	return c.SendStatus(fiber.StatusNoContent)
}

// Reset Password
// @Summary Reset Password
// @Tags Auth
// @Accept json
// @Produce json
// @Param token query string true "Token"
// @Param body body dto.ResetPasswordRequest true "Body"
// @Success 204
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	token := c.Query("token")
	var req dto.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil { return c.SendStatus(fiber.StatusBadRequest) }
	
	if err := h.Service.ResetPassword(token, req.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Send Verification Email
// @Summary Send Verification Email
// @Tags Auth
// @Security BearerAuth
// @Success 204
// @Router /auth/send-verification-email [post]
func (h *AuthHandler) SendVerificationEmail(c *fiber.Ctx) error {
	userID := c.Locals("userId").(uuid.UUID)
	h.Service.SendVerificationEmail(userID.String())
	return c.SendStatus(fiber.StatusNoContent)
}

// Verify Email
// @Summary Verify Email
// @Tags Auth
// @Param token query string true "Token"
// @Success 204
// @Router /auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if err := h.Service.VerifyEmail(token); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}