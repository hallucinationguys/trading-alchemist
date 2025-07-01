package dto

import (
	"time"

	"trading-alchemist/internal/domain/entities"

	"github.com/google/uuid"
)

// UserResponse represents a user in API responses
//
// swagger:model UserResponse
type UserResponse struct {
	// User ID
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID uuid.UUID `json:"id"`
	// Email address
	// example: user@example.com
	Email string `json:"email"`
	// Whether email is verified
	// example: true
	EmailVerified bool `json:"email_verified"`
	// First name
	// example: John
	FirstName *string `json:"first_name"`
	// Last name
	// example: Doe
	LastName *string `json:"last_name"`
	// Avatar URL
	// example: https://example.com/avatar.jpg
	AvatarURL *string `json:"avatar_url"`
	// Full name (computed from first and last name)
	// example: John Doe
	FullName string `json:"full_name"`
	// Display name
	// example: John Doe
	DisplayName string `json:"display_name"`
	// Whether user is active
	// example: true
	IsActive bool `json:"is_active"`
	// Account creation timestamp
	// example: 2023-12-01T10:00:00Z
	CreatedAt time.Time `json:"created_at"`
	// Last update timestamp
	// example: 2023-12-01T10:00:00Z
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
}

// UpdateUserRequest represents the request to update user information
//
// swagger:model UpdateUserRequest
type UpdateUserRequest struct {
	// First name (1-100 characters)
	// example: John
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	// Last name (1-100 characters)
	// example: Doe
	LastName *string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
	// Avatar URL
	// example: https://example.com/avatar.jpg
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
}

// GetUserResponse represents the response for getting user information
//
// swagger:model GetUserResponse
type GetUserResponse struct {
	// User information
	User UserResponse `json:"user"`
}

// UpdateUserResponse represents the response for updating user information
//
// swagger:model UpdateUserResponse
type UpdateUserResponse struct {
	// Updated user information
	User UserResponse `json:"user"`
}

// ToUserResponse converts a domain User entity to UserResponse DTO
func ToUserResponse(user *entities.User) UserResponse {
	return UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		AvatarURL:     user.AvatarURL,
		FullName:      user.FullName(),
		DisplayName:   user.DisplayName(),
		IsActive:      user.IsActive,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

// ToUserEntity converts a CreateUserRequest DTO to User entity
func (req *CreateUserRequest) ToUserEntity() *entities.User {
	return &entities.User{
		ID:            uuid.New(),
		Email:         req.Email,
		EmailVerified: false,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
} 