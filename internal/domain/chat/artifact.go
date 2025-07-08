package chat

import (
	"time"
	"trading-alchemist/internal/domain/shared"

	"github.com/google/uuid"
)

// Artifact represents generated content like code, documents, etc.
type Artifact struct {
	ID           uuid.UUID    `json:"id" db:"id"`
	MessageID    uuid.UUID    `json:"message_id" db:"message_id"`
	Title        string       `json:"title" db:"title"`
	Type         shared.ArtifactType `json:"type" db:"type"`
	Language     *string      `json:"language" db:"language"` // For code artifacts
	Content      string       `json:"content" db:"content"`
	ContentHash  string       `json:"content_hash" db:"content_hash"` // For deduplication
	Size         int64        `json:"size" db:"size"`                 // Content size in bytes
	IsPublic     bool         `json:"is_public" db:"is_public"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
} 