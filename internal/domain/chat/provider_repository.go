package chat

import (
	"context"

	"github.com/google/uuid"
)

// ProviderRepository defines the domain's interface for provider data operations.
type ProviderRepository interface {
	GetAll(ctx context.Context) ([]*Provider, error)
	GetActive(ctx context.Context) ([]*Provider, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Provider, error)
	GetByName(ctx context.Context, name string) (*Provider, error)
	Create(ctx context.Context, provider *Provider) (*Provider, error)
	// GetAvailableModelsForUser returns models with API key status for a user (optimized with JOINs)
	GetAvailableModelsForUser(ctx context.Context, userID uuid.UUID) ([]*Provider, error)
} 