package repositories

import (
	"context"

	"github.com/google/uuid"

	"trading-alchemist/internal/domain/entities"
)

type ModelRepository interface {
	GetActiveModelsByProviderID(ctx context.Context, providerID uuid.UUID) ([]*entities.Model, error)
	CreateModel(ctx context.Context, model *entities.Model) (*entities.Model, error)
	GetModelByName(ctx context.Context, providerID uuid.UUID, name string) (*entities.Model, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Model, error)
} 