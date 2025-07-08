package auth

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) (*User, error)
	
	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	
	// GetByEmail retrieves a user by their email
	GetByEmail(ctx context.Context, email string) (*User, error)
	
	// Update updates user information
	Update(ctx context.Context, user *User) (*User, error)
	
	// VerifyEmail marks a user's email as verified
	VerifyEmail(ctx context.Context, userID uuid.UUID) (*User, error)
	
	// Deactivate deactivates a user account
	Deactivate(ctx context.Context, userID uuid.UUID) error
} 