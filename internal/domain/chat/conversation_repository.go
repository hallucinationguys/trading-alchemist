package chat

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ConversationRepository interface {
	Create(ctx context.Context, conversation *Conversation) (*Conversation, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Conversation, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Conversation, error)
	Update(ctx context.Context, conversation *Conversation) (*Conversation, error)
	UpdateLastMessageAt(ctx context.Context, id uuid.UUID, lastMessageAt time.Time) error
	UpdateTitle(ctx context.Context, id uuid.UUID, title string) error
	Archive(ctx context.Context, id uuid.UUID) error
} 