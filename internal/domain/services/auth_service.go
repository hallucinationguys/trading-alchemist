package services

import (
	"context"

	"trading-alchemist/internal/domain/entities"

	"github.com/google/uuid"
)

type AuthService interface {
	// SendMagicLink generates and sends a magic link to the user's email
	SendMagicLink(ctx context.Context, email string, purpose entities.MagicLinkPurpose, ipAddress, userAgent string) error
	
	// VerifyMagicLink verifies a magic link and returns the user if valid
	VerifyMagicLink(ctx context.Context, token string) (*entities.User, error)
	
	// GenerateJWT generates a JWT token for the authenticated user
	GenerateJWT(ctx context.Context, user *entities.User) (string, error)
	
	// ValidateJWT validates a JWT token and returns the user ID
	ValidateJWT(ctx context.Context, token string) (uuid.UUID, error)
	
	// InvalidateUserSessions invalidates all sessions for a user
	InvalidateUserSessions(ctx context.Context, userID uuid.UUID) error
} 