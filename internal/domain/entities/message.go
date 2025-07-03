package entities

import (
	"time"

	"github.com/google/uuid"
)

// Message represents individual messages in conversations
type Message struct {
	ID             uuid.UUID     `json:"id" db:"id"`
	ConversationID uuid.UUID     `json:"conversation_id" db:"conversation_id"`
	ParentID       *uuid.UUID    `json:"parent_id" db:"parent_id"` // For message threads/edits
	Role           MessageRole   `json:"role" db:"role"`
	Content        string        `json:"content" db:"content"`
	ModelID        *uuid.UUID    `json:"model_id" db:"model_id"`      // Model used for this message
	TokenCount     *int          `json:"token_count" db:"token_count"`
	Cost           *float64      `json:"cost" db:"cost"` // Cost for this message
	Metadata       JSONB         `json:"metadata" db:"metadata"`       // Function calls, tool use, etc.
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
} 