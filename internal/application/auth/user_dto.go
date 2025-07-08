package auth

import (
	"time"
	"trading-alchemist/internal/domain/auth"

	"github.com/google/uuid"
)

// UserResponse represents a user in API responses
type UserResponse struct {
	// User ID
	ID uuid.UUID `json:"id"`
	// Email address
	Email string `json:"email"`
	// Whether email is verified
	EmailVerified bool `json:"email_verified"`
	// First name
	FirstName *string `json:"first_name"`
	// Last name
	LastName *string `json:"last_name"`
	// Avatar URL
	AvatarURL *string `json:"avatar_url"`
	// Full name (computed from first and last name)
	FullName string `json:"full_name"`
	// Display name
	DisplayName string `json:"display_name"`
	// Whether user is active
	IsActive bool `json:"is_active"`
	// Account creation timestamp
	CreatedAt time.Time `json:"created_at"`
	// Last update timestamp
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
}

// UpdateUserRequest represents the request to update user information
type UpdateUserRequest struct {
	// First name (1-100 characters)
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	// Last name (1-100 characters)
	LastName *string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
	// Avatar URL
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
}

// GetUserResponse represents the response for getting user information
type GetUserResponse struct {
	// User information
	User UserResponse `json:"user"`
}

// UpdateUserResponse represents the response for updating user information
type UpdateUserResponse struct {
	// Updated user information
	User UserResponse `json:"user"`
}

// ToUserResponse converts a domain User entity to UserResponse DTO
func ToUserResponse(user *auth.User) UserResponse {
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
func (req *CreateUserRequest) ToUserEntity() *auth.User {
	return &auth.User{
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