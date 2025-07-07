package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"trading-alchemist/internal/domain/entities"
	"trading-alchemist/internal/domain/repositories"
	"trading-alchemist/internal/infrastructure/repositories/postgres/sqlc"
	"trading-alchemist/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// MessageRepository implements the domain's MessageRepository interface using PostgreSQL.
type MessageRepository struct {
	queries *sqlc.Queries
}

// NewMessageRepository creates a new postgres message repository.
func NewMessageRepository(db sqlc.DBTX) repositories.MessageRepository {
	return &MessageRepository{
		queries: sqlc.New(db),
	}
}

func (r *MessageRepository) Create(ctx context.Context, message *entities.Message) (*entities.Message, error) {
	params := sqlc.CreateMessageParams{
		ConversationID: pgtype.UUID{Bytes: message.ConversationID, Valid: true},
		Role:           string(message.Role),
		Content:        pgtype.Text{String: message.Content, Valid: true},
	}
	if message.ParentID != nil {
		params.ParentID = pgtype.UUID{Bytes: *message.ParentID, Valid: true}
	}
	if message.ModelID != nil {
		params.ModelID = pgtype.UUID{Bytes: *message.ModelID, Valid: true}
	}
	if message.TokenCount != nil {
		params.TokenCount = pgtype.Int4{Int32: int32(*message.TokenCount), Valid: true}
	}
	if message.Cost != nil {
		costNumeric := pgtype.Numeric{}
		err := costNumeric.Scan(*message.Cost)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message cost: %w", err)
		}
		params.Cost = costNumeric
	}
	if message.Metadata != nil {
		metadataJSON, err := json.Marshal(message.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal message metadata: %w", err)
		}
		params.Metadata = metadataJSON
	}

	sqlcMessage, err := r.queries.CreateMessage(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}
	return sqlcMessageToEntity(&sqlcMessage), nil
}

func (r *MessageRepository) GetByConversationID(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]*entities.Message, error) {
	params := sqlc.GetMessagesByConversationIDParams{
		ConversationID: pgtype.UUID{Bytes: conversationID, Valid: true},
		Limit:          int32(limit),
		Offset:         int32(offset),
	}
	sqlcMessages, err := r.queries.GetMessagesByConversationID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages by conversation ID: %w", err)
	}

	messages := make([]*entities.Message, len(sqlcMessages))
	for i, m := range sqlcMessages {
		messages[i] = sqlcMessageToEntity(&m)
	}
	return messages, nil
}

func (r *MessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Message, error) {
	messageUUID := pgtype.UUID{Bytes: id, Valid: true}
	sqlcMessage, err := r.queries.GetMessageByID(ctx, messageUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrMessageNotFound
		}
		return nil, fmt.Errorf("failed to get message by ID: %w", err)
	}
	return sqlcMessageToEntity(&sqlcMessage), nil
}

func (r *MessageRepository) GetThread(ctx context.Context, parentID uuid.UUID) ([]*entities.Message, error) {
	parentUUID := pgtype.UUID{Bytes: parentID, Valid: true}
	sqlcMessages, err := r.queries.GetMessageThread(ctx, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message thread: %w", err)
	}

	messages := make([]*entities.Message, len(sqlcMessages))
	for i, m := range sqlcMessages {
		messages[i] = sqlcMessageToEntity(&m)
	}
	return messages, nil
}

func (r *MessageRepository) Update(ctx context.Context, message *entities.Message) (*entities.Message, error) {
	params := sqlc.UpdateMessageParams{
		ID: pgtype.UUID{Bytes: message.ID, Valid: true},
	}
	// Only update fields that are not nil in the entity
	if message.Content != "" {
		params.Content = pgtype.Text{String: message.Content, Valid: true}
	}
	if message.TokenCount != nil {
		params.TokenCount = pgtype.Int4{Int32: int32(*message.TokenCount), Valid: true}
	}
	if message.Cost != nil {
		costNumeric := pgtype.Numeric{}
		err := costNumeric.Scan(*message.Cost)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message cost: %w", err)
		}
		params.Cost = costNumeric
	}
	if message.Metadata != nil {
		metadataJSON, err := json.Marshal(message.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal message metadata: %w", err)
		}
		params.Metadata = metadataJSON
	}

	sqlcMessage, err := r.queries.UpdateMessage(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrMessageNotFound
		}
		return nil, fmt.Errorf("failed to update message: %w", err)
	}
	return sqlcMessageToEntity(&sqlcMessage), nil
}

func (r *MessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	messageUUID := pgtype.UUID{Bytes: id, Valid: true}
	return r.queries.DeleteMessage(ctx, messageUUID)
}

func (r *MessageRepository) GetByConversationIDWithCursor(ctx context.Context, conversationID uuid.UUID, cursor *time.Time, limit int) ([]*entities.Message, error) {
	params := sqlc.GetMessagesByConversationIDWithCursorParams{
		ConversationID: pgtype.UUID{Bytes: conversationID, Valid: true},
		Limit:          int32(limit),
	}
	if cursor != nil {
		params.CreatedAt = pgtype.Timestamptz{Time: *cursor, Valid: true}
	}

	sqlcMessages, err := r.queries.GetMessagesByConversationIDWithCursor(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages by conversation ID with cursor: %w", err)
	}

	messages := make([]*entities.Message, len(sqlcMessages))
	for i, m := range sqlcMessages {
		messages[i] = sqlcMessageToEntity(&m)
	}
	return messages, nil
}

func (r *MessageRepository) CountByConversationID(ctx context.Context, conversationID uuid.UUID) (int, error) {
	convUUID := pgtype.UUID{Bytes: conversationID, Valid: true}
	count, err := r.queries.CountMessagesByConversationID(ctx, convUUID)
	if err != nil {
		return 0, fmt.Errorf("failed to count messages by conversation ID: %w", err)
	}
	return int(count), nil
}

func sqlcMessageToEntity(m *sqlc.Message) *entities.Message {
	msg := &entities.Message{
		Role:      entities.MessageRole(m.Role),
		Content:   m.Content.String,
		CreatedAt: m.CreatedAt.Time,
		UpdatedAt: m.UpdatedAt.Time,
	}

	if m.ID.Valid {
		msg.ID = m.ID.Bytes
	}
	if m.ConversationID.Valid {
		msg.ConversationID = m.ConversationID.Bytes
	}
	if m.ParentID.Valid {
		parentID := m.ParentID.Bytes
		msg.ParentID = (*uuid.UUID)(&parentID)
	}
	if m.ModelID.Valid {
		modelID := m.ModelID.Bytes
		msg.ModelID = (*uuid.UUID)(&modelID)
	}
	if m.TokenCount.Valid {
		tokenCount := int(m.TokenCount.Int32)
		msg.TokenCount = &tokenCount
	}
	if m.Cost.Valid {
		var cost float64
		if err := m.Cost.Scan(&cost); err == nil {
			msg.Cost = &cost
		}
	}
	if m.Metadata != nil {
		var metadata entities.JSONB
		if err := json.Unmarshal(m.Metadata, &metadata); err == nil {
			msg.Metadata = metadata
		}
	}
	return msg
} 