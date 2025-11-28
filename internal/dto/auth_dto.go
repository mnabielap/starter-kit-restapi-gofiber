package dto

// RegisterRequest
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,alphanum"`
}

// LoginRequest
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LogoutRequest
type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// RefreshTokenRequest
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// ForgotPasswordRequest
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest
type ResetPasswordRequest struct {
	Token    string `query:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8,alphanum"`
}

// VerifyEmailRequest
type VerifyEmailRequest struct {
	Token string `query:"token" validate:"required"`
}

// AuthResponse
type AuthResponse struct {
	User   UserResponse   `json:"user"`
	Tokens TokensResponse `json:"tokens"`
}

type TokensResponse struct {
	Access  TokenDetails `json:"access"`
	Refresh TokenDetails `json:"refresh"`
}

type TokenDetails struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}