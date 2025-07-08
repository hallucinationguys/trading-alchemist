package postgres

import (
	"context"

	"trading-alchemist/internal/domain/chat"
	"trading-alchemist/internal/infrastructure/repositories/postgres/shared/sqlc"
	"trading-alchemist/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ModelRepository struct {
	q *sqlc.Queries
}

func NewModelRepository(db sqlc.DBTX) *ModelRepository {
	return &ModelRepository{q: sqlc.New(db)}
}

func (r *ModelRepository) GetActiveModelsByProviderID(ctx context.Context, providerID uuid.UUID) ([]*chat.Model, error) {
	pgProviderID := pgtype.UUID{Bytes: providerID, Valid: true}
	dbModels, err := r.q.GetActiveModelsByProviderID(ctx, pgProviderID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []*chat.Model{}, nil
		}
		return nil, err
	}

	models := make([]*chat.Model, len(dbModels))
	for i, dbModel := range dbModels {
		models[i] = &chat.Model{
			ID:                dbModel.ID.Bytes,
			ProviderID:        dbModel.ProviderID.Bytes,
			Name:              dbModel.Name,
			DisplayName:       dbModel.DisplayName,
			SupportsFunctions: dbModel.SupportsFunctions.Bool,
			SupportsVision:    dbModel.SupportsVision.Bool,
			IsActive:          dbModel.IsActive.Bool,
			CreatedAt:         dbModel.CreatedAt.Time,
			UpdatedAt:         dbModel.UpdatedAt.Time,
		}
	}
	return models, nil
}

func (r *ModelRepository) CreateModel(ctx context.Context, model *chat.Model) (*chat.Model, error) {
	dbModel, err := r.q.CreateModel(ctx, sqlc.CreateModelParams{
		ProviderID:        pgtype.UUID{Bytes: model.ProviderID, Valid: true},
		Name:              model.Name,
		DisplayName:       model.DisplayName,
		SupportsFunctions: pgtype.Bool{Bool: model.SupportsFunctions, Valid: true},
		SupportsVision:    pgtype.Bool{Bool: model.SupportsVision, Valid: true},
		IsActive:          pgtype.Bool{Bool: model.IsActive, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &chat.Model{
		ID:                dbModel.ID.Bytes,
		ProviderID:        dbModel.ProviderID.Bytes,
		Name:              dbModel.Name,
		DisplayName:       dbModel.DisplayName,
		SupportsFunctions: dbModel.SupportsFunctions.Bool,
		SupportsVision:    dbModel.SupportsVision.Bool,
		IsActive:          dbModel.IsActive.Bool,
		CreatedAt:         dbModel.CreatedAt.Time,
		UpdatedAt:         dbModel.UpdatedAt.Time,
	}, nil
}

func (r *ModelRepository) GetModelByName(ctx context.Context, providerID uuid.UUID, name string) (*chat.Model, error) {
	dbModel, err := r.q.GetModelByName(ctx, sqlc.GetModelByNameParams{
		ProviderID: pgtype.UUID{Bytes: providerID, Valid: true},
		Name:       name,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrModelNotFound
		}
		return nil, err
	}

	return &chat.Model{
		ID:                dbModel.ID.Bytes,
		ProviderID:        dbModel.ProviderID.Bytes,
		Name:              dbModel.Name,
		DisplayName:       dbModel.DisplayName,
		SupportsFunctions: dbModel.SupportsFunctions.Bool,
		SupportsVision:    dbModel.SupportsVision.Bool,
		IsActive:          dbModel.IsActive.Bool,
		CreatedAt:         dbModel.CreatedAt.Time,
		UpdatedAt:         dbModel.UpdatedAt.Time,
	}, nil
}

func (r *ModelRepository) GetByID(ctx context.Context, id uuid.UUID) (*chat.Model, error) {
	dbModel, err := r.q.GetModelByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrModelNotFound
		}
		return nil, err
	}

	return &chat.Model{
		ID:                dbModel.ID.Bytes,
		ProviderID:        dbModel.ProviderID.Bytes,
		Name:              dbModel.Name,
		DisplayName:       dbModel.DisplayName,
		SupportsFunctions: dbModel.SupportsFunctions.Bool,
		SupportsVision:    dbModel.SupportsVision.Bool,
		IsActive:          dbModel.IsActive.Bool,
		CreatedAt:         dbModel.CreatedAt.Time,
		UpdatedAt:         dbModel.UpdatedAt.Time,
	}, nil
} 