package repositories

import (
	"context"
	"time"
	"trading-alchemist/internal/domain/entities"

	"github.com/google/uuid"
)

type MessageRepository interface {
	Create(ctx context.Context, message *entities.Message) (*entities.Message, error)
	GetByConversationID(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]*entities.Message, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Message, error)
	GetThread(ctx context.Context, parentID uuid.UUID) ([]*entities.Message, error)
	Update(ctx context.Context, message *entities.Message) (*entities.Message, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// For large conversations - get paginated with cursor
	GetByConversationIDWithCursor(ctx context.Context, conversationID uuid.UUID, cursor *time.Time, limit int) ([]*entities.Message, error)
} 