package auth

import (
	"context"

	"github.com/google/uuid"
)

type MagicLinkRepository interface {
	// Create creates a new magic link
	Create(ctx context.Context, magicLink *MagicLink) (*MagicLink, error)
	
	// GetByToken retrieves a magic link by its token (with user info)
	GetByToken(ctx context.Context, token string) (*MagicLink, *User, error)
	
	// MarkAsUsed marks a magic link as used
	MarkAsUsed(ctx context.Context, linkID uuid.UUID) (*MagicLink, error)
	
	// InvalidateUserLinks invalidates all unused magic links for a user with specific purpose
	InvalidateUserLinks(ctx context.Context, userID uuid.UUID, purpose MagicLinkPurpose) error
	
	// CleanupExpired removes expired magic links
	CleanupExpired(ctx context.Context) error
} 