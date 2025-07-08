package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"trading-alchemist/internal/domain/chat"
	"trading-alchemist/internal/domain/shared"
	"trading-alchemist/internal/infrastructure/repositories/postgres/shared/sqlc"
	"trading-alchemist/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ConversationRepository implements the domain's ConversationRepository interface using PostgreSQL.
type ConversationRepository struct {
	queries *sqlc.Queries
}

// NewConversationRepository creates a new postgres conversation repository.
func NewConversationRepository(db sqlc.DBTX) chat.ConversationRepository {
	return &ConversationRepository{
		queries: sqlc.New(db),
	}
}

func (r *ConversationRepository) Create(ctx context.Context, conversation *chat.Conversation) (*chat.Conversation, error) {
	params := sqlc.CreateConversationParams{
		UserID:  pgtype.UUID{Bytes: conversation.UserID, Valid: true},
		Title:   conversation.Title,
		ModelID: pgtype.UUID{Bytes: conversation.ModelID, Valid: true},
	}
	if conversation.SystemPrompt != nil {
		params.SystemPrompt = pgtype.Text{String: *conversation.SystemPrompt, Valid: true}
	}
	if conversation.Settings != nil {
		settingsJSON, err := json.Marshal(conversation.Settings)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal conversation settings: %w", err)
		}
		params.Settings = settingsJSON
	}

	sqlcConv, err := r.queries.CreateConversation(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return sqlcConversationToEntity(&sqlcConv), nil
}

func (r *ConversationRepository) GetByID(ctx context.Context, id uuid.UUID) (*chat.Conversation, error) {
	convUUID := pgtype.UUID{Bytes: id, Valid: true}
	sqlcConv, err := r.queries.GetConversationByID(ctx, convUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrConversationNotFound
		}
		return nil, fmt.Errorf("failed to get conversation by ID: %w", err)
	}
	return sqlcConversationToEntity(&sqlcConv), nil
}

func (r *ConversationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*chat.Conversation, error) {
	params := sqlc.GetConversationsByUserIDParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	sqlcConvs, err := r.queries.GetConversationsByUserID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations by user ID: %w", err)
	}

	convs := make([]*chat.Conversation, len(sqlcConvs))
	for i, c := range sqlcConvs {
		convs[i] = sqlcConversationToEntity(&c)
	}
	return convs, nil
}

func (r *ConversationRepository) Update(ctx context.Context, conversation *chat.Conversation) (*chat.Conversation, error) {
	params := sqlc.UpdateConversationParams{
		ID:      pgtype.UUID{Bytes: conversation.ID, Valid: true},
		Title:   conversation.Title,
		ModelID: pgtype.UUID{Bytes: conversation.ModelID, Valid: true},
	}
	if conversation.SystemPrompt != nil {
		params.SystemPrompt = pgtype.Text{String: *conversation.SystemPrompt, Valid: true}
	}
	if conversation.Settings != nil {
		settingsJSON, err := json.Marshal(conversation.Settings)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal conversation settings: %w", err)
		}
		params.Settings = settingsJSON
	}

	sqlcConv, err := r.queries.UpdateConversation(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrConversationNotFound
		}
		return nil, fmt.Errorf("failed to update conversation: %w", err)
	}

	return sqlcConversationToEntity(&sqlcConv), nil
}

func (r *ConversationRepository) UpdateLastMessageAt(ctx context.Context, id uuid.UUID, lastMessageAt time.Time) error {
	params := sqlc.UpdateConversationLastMessageAtParams{
		ID:            pgtype.UUID{Bytes: id, Valid: true},
		LastMessageAt: pgtype.Timestamptz{Time: lastMessageAt, Valid: true},
	}
	return r.queries.UpdateConversationLastMessageAt(ctx, params)
}

func (r *ConversationRepository) UpdateTitle(ctx context.Context, id uuid.UUID, title string) error {
	params := sqlc.UpdateConversationTitleParams{
		ID:    pgtype.UUID{Bytes: id, Valid: true},
		Title: title,
	}
	return r.queries.UpdateConversationTitle(ctx, params)
}

func (r *ConversationRepository) Archive(ctx context.Context, id uuid.UUID) error {
	convUUID := pgtype.UUID{Bytes: id, Valid: true}
	return r.queries.ArchiveConversation(ctx, convUUID)
}

func sqlcConversationToEntity(c *sqlc.Conversation) *chat.Conversation {
	conv := &chat.Conversation{
		Title:      c.Title,
		IsArchived: c.IsArchived.Bool,
		CreatedAt:  c.CreatedAt.Time,
		UpdatedAt:  c.UpdatedAt.Time,
	}

	if c.ID.Valid {
		conv.ID = c.ID.Bytes
	}
	if c.UserID.Valid {
		conv.UserID = c.UserID.Bytes
	}
	if c.ModelID.Valid {
		conv.ModelID = c.ModelID.Bytes
	}
	if c.SystemPrompt.Valid {
		conv.SystemPrompt = &c.SystemPrompt.String
	}
	if c.Settings != nil {
		var settings shared.JSONB
		if err := json.Unmarshal(c.Settings, &settings); err == nil {
			conv.Settings = settings
		}
	}
	if c.LastMessageAt.Valid {
		conv.LastMessageAt = &c.LastMessageAt.Time
	}

	return conv
} 