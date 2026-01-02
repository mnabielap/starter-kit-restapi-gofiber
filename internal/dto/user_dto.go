package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Role            string    `json:"role"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	CreatedAt       time.Time `json:"createdAt"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=user admin"`
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
}

// UserQueryParams holds search, filter, and sort options
type UserQueryParams struct {
	Search string
	Scope  string
	Role   string
	SortBy string
	Limit  int
	Page   int
}