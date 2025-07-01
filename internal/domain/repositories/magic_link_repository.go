package repositories

import (
	"context"

	"trading-alchemist/internal/domain/entities"

	"github.com/google/uuid"
)

type MagicLinkRepository interface {
	// Create creates a new magic link
	Create(ctx context.Context, magicLink *entities.MagicLink) (*entities.MagicLink, error)
	
	// GetByToken retrieves a magic link by its token (with user info)
	GetByToken(ctx context.Context, token string) (*entities.MagicLink, *entities.User, error)
	
	// MarkAsUsed marks a magic link as used
	MarkAsUsed(ctx context.Context, linkID uuid.UUID) (*entities.MagicLink, error)
	
	// InvalidateUserLinks invalidates all unused magic links for a user with specific purpose
	InvalidateUserLinks(ctx context.Context, userID uuid.UUID, purpose entities.MagicLinkPurpose) error
	
	// CleanupExpired removes expired magic links
	CleanupExpired(ctx context.Context) error
} 