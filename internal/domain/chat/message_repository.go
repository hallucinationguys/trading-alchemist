package chat

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MessageRepository interface {
	Create(ctx context.Context, message *Message) (*Message, error)
	GetByConversationID(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]*Message, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Message, error)
	GetThread(ctx context.Context, parentID uuid.UUID) ([]*Message, error)
	Update(ctx context.Context, message *Message) (*Message, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// For large conversations - get paginated with cursor
	GetByConversationIDWithCursor(ctx context.Context, conversationID uuid.UUID, cursor *time.Time, limit int) ([]*Message, error)
	// Count messages in a conversation for title generation
	CountByConversationID(ctx context.Context, conversationID uuid.UUID) (int, error)
} 