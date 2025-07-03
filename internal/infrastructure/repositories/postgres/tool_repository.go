package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"trading-alchemist/internal/domain/entities"
	"trading-alchemist/internal/domain/repositories"
	"trading-alchemist/internal/infrastructure/repositories/postgres/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ToolRepository implements the domain's ToolRepository interface using PostgreSQL.
type ToolRepository struct {
	queries *sqlc.Queries
}

// NewToolRepository creates a new postgres tool repository.
func NewToolRepository(db sqlc.DBTX) repositories.ToolRepository {
	return &ToolRepository{
		queries: sqlc.New(db),
	}
}

func (r *ToolRepository) GetAvailableTools(ctx context.Context, providerID *uuid.UUID) ([]*entities.Tool, error) {
	var providerUUID pgtype.UUID
	if providerID != nil {
		providerUUID = pgtype.UUID{Bytes: *providerID, Valid: true}
	}

	sqlcTools, err := r.queries.GetAvailableTools(ctx, providerUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available tools: %w", err)
	}

	tools := make([]*entities.Tool, len(sqlcTools))
	for i, t := range sqlcTools {
		tools[i] = sqlcToolToEntity(&t)
	}
	return tools, nil
}

func (r *ToolRepository) LogToolUsage(ctx context.Context, messageTool *entities.MessageTool) error {
	inputJSON, err := json.Marshal(messageTool.Input)
	if err != nil {
		return fmt.Errorf("failed to marshal tool input: %w", err)
	}
	outputJSON, err := json.Marshal(messageTool.Output)
	if err != nil {
		return fmt.Errorf("failed to marshal tool output: %w", err)
	}

	params := sqlc.LogToolUsageParams{
		MessageID:  pgtype.UUID{Bytes: messageTool.MessageID, Valid: true},
		ToolID:     pgtype.UUID{Bytes: messageTool.ToolID, Valid: true},
		Input:      inputJSON,
		Output:     outputJSON,
		ExecutedAt: pgtype.Timestamptz{Time: messageTool.ExecutedAt, Valid: true},
		Duration:   pgtype.Int8{Int64: messageTool.Duration, Valid: true},
		Success:    pgtype.Bool{Bool: messageTool.Success, Valid: true},
	}
	if messageTool.Error != nil {
		params.Error = pgtype.Text{String: *messageTool.Error, Valid: true}
	}

	_, err = r.queries.LogToolUsage(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to log tool usage: %w", err)
	}
	return nil
}

func sqlcToolToEntity(t *sqlc.Tool) *entities.Tool {
	tool := &entities.Tool{
		Name:        t.Name,
		Description: t.Description.String,
		IsActive:    t.IsActive.Bool,
		CreatedAt:   t.CreatedAt.Time,
		UpdatedAt:   t.UpdatedAt.Time,
	}
	if t.Schema != nil {
		var schema entities.JSONB
		if err := json.Unmarshal(t.Schema, &schema); err == nil {
			tool.Schema = schema
		}
	}

	if t.ID.Valid {
		tool.ID = t.ID.Bytes
	}
	if t.ProviderID.Valid {
		providerID := t.ProviderID.Bytes
		tool.ProviderID = (*uuid.UUID)(&providerID)
	}

	return tool
}

// TODO: Implement methods using sqlc generated queries. 