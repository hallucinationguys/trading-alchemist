package chat

import (
	"context"

	"github.com/google/uuid"
)

type ArtifactRepository interface {
	Create(ctx context.Context, artifact *Artifact) (*Artifact, error)
	GetByMessageID(ctx context.Context, messageID uuid.UUID) ([]*Artifact, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Artifact, error)
	Update(ctx context.Context, artifact *Artifact) (*Artifact, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetPublicArtifacts(ctx context.Context, limit, offset int) ([]*Artifact, error)
} 