package postgres

import (
	"context"
	"fmt"

	"trading-alchemist/internal/domain/chat"
	"trading-alchemist/internal/infrastructure/repositories/postgres/shared/sqlc"
	"trading-alchemist/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ProviderRepository implements the domain's ProviderRepository interface using PostgreSQL.
type ProviderRepository struct {
	queries *sqlc.Queries
}

// NewProviderRepository creates a new postgres provider repository.
func NewProviderRepository(db sqlc.DBTX) chat.ProviderRepository {
	return &ProviderRepository{
		queries: sqlc.New(db),
	}
}

func (r *ProviderRepository) GetAll(ctx context.Context) ([]*chat.Provider, error) {
	sqlcProviders, err := r.queries.GetAllProviders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all providers: %w", err)
	}

	providers := make([]*chat.Provider, len(sqlcProviders))
	for i, p := range sqlcProviders {
		providers[i] = sqlcProviderToEntity(&p)
	}

	return providers, nil
}

func (r *ProviderRepository) GetActive(ctx context.Context) ([]*chat.Provider, error) {
	sqlcProviders, err := r.queries.GetActiveProviders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active providers: %w", err)
	}

	providers := make([]*chat.Provider, len(sqlcProviders))
	for i, p := range sqlcProviders {
		providers[i] = sqlcProviderToEntity(&p)
	}

	return providers, nil
}

func (r *ProviderRepository) GetByID(ctx context.Context, id uuid.UUID) (*chat.Provider, error) {
	providerUUID := pgtype.UUID{Bytes: id, Valid: true}
	sqlcProvider, err := r.queries.GetProviderByID(ctx, providerUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrProviderNotFound
		}
		return nil, fmt.Errorf("failed to get provider by ID: %w", err)
	}
	return sqlcProviderToEntity(&sqlcProvider), nil
}

func (r *ProviderRepository) GetByName(ctx context.Context, name string) (*chat.Provider, error) {
	sqlcProvider, err := r.queries.GetProviderByName(ctx, name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrProviderNotFound
		}
		return nil, fmt.Errorf("failed to get provider by name: %w", err)
	}
	return sqlcProviderToEntity(&sqlcProvider), nil
}

func (r *ProviderRepository) Create(ctx context.Context, provider *chat.Provider) (*chat.Provider, error) {
	params := sqlc.CreateProviderParams{
		Name:        provider.Name,
		DisplayName: provider.DisplayName,
		IsActive:    pgtype.Bool{Bool: provider.IsActive, Valid: true},
	}
	sqlcProvider, err := r.queries.CreateProvider(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}
	return sqlcProviderToEntity(&sqlcProvider), nil
}

// GetAvailableModelsForUser returns providers with models and API key status for a user
func (r *ProviderRepository) GetAvailableModelsForUser(ctx context.Context, userID uuid.UUID) ([]*chat.Provider, error) {
	userIDPg := pgtype.UUID{}
	if err := userIDPg.Scan(userID.String()); err != nil {
		return nil, fmt.Errorf("failed to convert user ID: %w", err)
	}
	
	rows, err := r.queries.GetAvailableModelsForUser(ctx, userIDPg)
	if err != nil {
		return nil, fmt.Errorf("failed to get available models for user: %w", err)
	}
	
	// Transform flat JOIN result into hierarchical structure
	providerMap := make(map[string]*chat.Provider)
	
	for _, row := range rows {
		providerIDStr := row.ProviderID.String()
		
		// Get or create provider
		provider, exists := providerMap[providerIDStr]
		if !exists {
			provider = &chat.Provider{
				ID:          uuid.MustParse(row.ProviderID.String()),
				Name:        row.ProviderName,
				DisplayName: row.ProviderDisplayName,
				IsActive:    true, // Only active providers are returned
				Models:      []*chat.Model{},
			}
			providerMap[providerIDStr] = provider
		}
		
		// Add model if it exists (LEFT JOIN might have NULL models)
		if row.ModelID.Valid {
			model := &chat.Model{
				ID:                uuid.MustParse(row.ModelID.String()),
				ProviderID:        provider.ID,
				Name:             row.ModelName,
				DisplayName:      row.ModelDisplayName,
				SupportsFunctions: row.ModelSupportsFunctions.Bool,
				SupportsVision:   row.ModelSupportsVision.Bool,
				IsActive:         true, // Only active models are returned
			}
			
			provider.Models = append(provider.Models, model)
		}
	}
	
	// Convert map to slice
	result := make([]*chat.Provider, 0, len(providerMap))
	for _, provider := range providerMap {
		result = append(result, provider)
	}
	
	return result, nil
}

// sqlcProviderToEntity converts a SQLC Provider to a domain Provider entity.
func sqlcProviderToEntity(p *sqlc.Provider) *chat.Provider {
	provider := &chat.Provider{
		Name:        p.Name,
		DisplayName: p.DisplayName,
		IsActive:    p.IsActive.Bool,
		CreatedAt:   p.CreatedAt.Time,
		UpdatedAt:   p.UpdatedAt.Time,
	}

	if p.ID.Valid {
		provider.ID = p.ID.Bytes
	}

	return provider
}