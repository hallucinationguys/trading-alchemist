package postgres

import (
	"context"
	"fmt"

	"trading-alchemist/internal/domain/entities"
	"trading-alchemist/internal/domain/repositories"
	"trading-alchemist/internal/infrastructure/repositories/postgres/sqlc"
	"trading-alchemist/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ArtifactRepository implements the domain's ArtifactRepository interface using PostgreSQL.
type ArtifactRepository struct {
	queries *sqlc.Queries
}

// NewArtifactRepository creates a new postgres artifact repository.
func NewArtifactRepository(db sqlc.DBTX) repositories.ArtifactRepository {
	return &ArtifactRepository{
		queries: sqlc.New(db),
	}
}

func (r *ArtifactRepository) Create(ctx context.Context, artifact *entities.Artifact) (*entities.Artifact, error) {
	params := sqlc.CreateArtifactParams{
		MessageID:   pgtype.UUID{Bytes: artifact.MessageID, Valid: true},
		Title:       artifact.Title,
		Type:        string(artifact.Type),
		Content:     pgtype.Text{String: artifact.Content, Valid: true},
		ContentHash: pgtype.Text{String: artifact.ContentHash, Valid: true},
		Size:        pgtype.Int8{Int64: artifact.Size, Valid: true},
		IsPublic:    pgtype.Bool{Bool: artifact.IsPublic, Valid: true},
	}
	if artifact.Language != nil {
		params.Language = pgtype.Text{String: *artifact.Language, Valid: true}
	}

	sqlcArtifact, err := r.queries.CreateArtifact(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create artifact: %w", err)
	}
	return sqlcArtifactToEntity(&sqlcArtifact), nil
}

func (r *ArtifactRepository) GetByMessageID(ctx context.Context, messageID uuid.UUID) ([]*entities.Artifact, error) {
	messageUUID := pgtype.UUID{Bytes: messageID, Valid: true}
	sqlcArtifacts, err := r.queries.GetArtifactsByMessageID(ctx, messageUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get artifacts by message ID: %w", err)
	}

	artifacts := make([]*entities.Artifact, len(sqlcArtifacts))
	for i, a := range sqlcArtifacts {
		artifacts[i] = sqlcArtifactToEntity(&a)
	}
	return artifacts, nil
}

func (r *ArtifactRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Artifact, error) {
	artifactUUID := pgtype.UUID{Bytes: id, Valid: true}
	sqlcArtifact, err := r.queries.GetArtifactByID(ctx, artifactUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrArtifactNotFound
		}
		return nil, fmt.Errorf("failed to get artifact by ID: %w", err)
	}
	return sqlcArtifactToEntity(&sqlcArtifact), nil
}

func (r *ArtifactRepository) Update(ctx context.Context, artifact *entities.Artifact) (*entities.Artifact, error) {
	params := sqlc.UpdateArtifactParams{
		ID:          pgtype.UUID{Bytes: artifact.ID, Valid: true},
		Title:       artifact.Title,
		Content:     pgtype.Text{String: artifact.Content, Valid: true},
		ContentHash: pgtype.Text{String: artifact.ContentHash, Valid: true},
		Size:        pgtype.Int8{Int64: artifact.Size, Valid: true},
		IsPublic:    pgtype.Bool{Bool: artifact.IsPublic, Valid: true},
	}

	sqlcArtifact, err := r.queries.UpdateArtifact(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrArtifactNotFound
		}
		return nil, fmt.Errorf("failed to update artifact: %w", err)
	}
	return sqlcArtifactToEntity(&sqlcArtifact), nil
}

func (r *ArtifactRepository) Delete(ctx context.Context, id uuid.UUID) error {
	artifactUUID := pgtype.UUID{Bytes: id, Valid: true}
	return r.queries.DeleteArtifact(ctx, artifactUUID)
}

func (r *ArtifactRepository) GetPublicArtifacts(ctx context.Context, limit, offset int) ([]*entities.Artifact, error) {
	params := sqlc.GetPublicArtifactsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	sqlcArtifacts, err := r.queries.GetPublicArtifacts(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get public artifacts: %w", err)
	}

	artifacts := make([]*entities.Artifact, len(sqlcArtifacts))
	for i, a := range sqlcArtifacts {
		artifacts[i] = sqlcArtifactToEntity(&a)
	}
	return artifacts, nil
}

func sqlcArtifactToEntity(a *sqlc.Artifact) *entities.Artifact {
	artifact := &entities.Artifact{
		Title:       a.Title,
		Type:        entities.ArtifactType(a.Type),
		Content:     a.Content.String,
		ContentHash: a.ContentHash.String,
		Size:        a.Size.Int64,
		IsPublic:    a.IsPublic.Bool,
		CreatedAt:   a.CreatedAt.Time,
		UpdatedAt:   a.UpdatedAt.Time,
	}

	if a.ID.Valid {
		artifact.ID = a.ID.Bytes
	}
	if a.MessageID.Valid {
		artifact.MessageID = a.MessageID.Bytes
	}
	if a.Language.Valid {
		artifact.Language = &a.Language.String
	}

	return artifact
} 