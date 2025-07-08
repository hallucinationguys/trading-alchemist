package chat

import (
	"context"

	"github.com/google/uuid"

)

type ModelRepository interface {
	GetActiveModelsByProviderID(ctx context.Context, providerID uuid.UUID) ([]*Model, error)
	CreateModel(ctx context.Context, model *Model) (*Model, error)
	GetModelByName(ctx context.Context, providerID uuid.UUID, name string) (*Model, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Model, error)
} 