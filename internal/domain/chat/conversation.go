package chat

import (
	"time"
	"trading-alchemist/internal/domain/shared"

	"github.com/google/uuid"
)

// Conversation represents a chat session
type Conversation struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	Title         string     `json:"title" db:"title"`
	ModelID       uuid.UUID  `json:"model_id" db:"model_id"` // Current model
	SystemPrompt  *string    `json:"system_prompt" db:"system_prompt"`
	Settings      shared.JSONB      `json:"settings" db:"settings"` // Temperature, max_tokens, etc.
	IsArchived    bool       `json:"is_archived" db:"is_archived"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	LastMessageAt *time.Time `json:"last_message_at" db:"last_message_at"`
} 