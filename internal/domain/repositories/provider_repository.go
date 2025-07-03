package repositories

import (
	"context"
	"trading-alchemist/internal/domain/entities"

	"github.com/google/uuid"
)

// ProviderRepository defines the domain's interface for provider data operations.
type ProviderRepository interface {
	GetAll(ctx context.Context) ([]*entities.Provider, error)
	GetActive(ctx context.Context) ([]*entities.Provider, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Provider, error)
	GetByName(ctx context.Context, name string) (*entities.Provider, error)
	Create(ctx context.Context, provider *entities.Provider) (*entities.Provider, error)
	// GetAvailableModelsForUser returns models with API key status for a user (optimized with JOINs)
	GetAvailableModelsForUser(ctx context.Context, userID uuid.UUID) ([]*entities.Provider, error)
} 