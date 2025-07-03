package repositories

import (
	"context"
	"time"
	"trading-alchemist/internal/domain/entities"

	"github.com/google/uuid"
)

type ConversationRepository interface {
	Create(ctx context.Context, conversation *entities.Conversation) (*entities.Conversation, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Conversation, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Conversation, error)
	Update(ctx context.Context, conversation *entities.Conversation) (*entities.Conversation, error)
	UpdateLastMessageAt(ctx context.Context, id uuid.UUID, lastMessageAt time.Time) error
	Archive(ctx context.Context, id uuid.UUID) error
} 