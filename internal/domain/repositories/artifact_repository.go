package repositories

import (
	"context"
	"trading-alchemist/internal/domain/entities"

	"github.com/google/uuid"
)

type ArtifactRepository interface {
	Create(ctx context.Context, artifact *entities.Artifact) (*entities.Artifact, error)
	GetByMessageID(ctx context.Context, messageID uuid.UUID) ([]*entities.Artifact, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Artifact, error)
	Update(ctx context.Context, artifact *entities.Artifact) (*entities.Artifact, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetPublicArtifacts(ctx context.Context, limit, offset int) ([]*entities.Artifact, error)
} 