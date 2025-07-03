package dto

import (
	"time"

	"github.com/google/uuid"
)

// --- Request DTOs ---

// CreateConversationRequest represents the request to create a new conversation.
type CreateConversationRequest struct {
	Title     string    `json:"title" validate:"required,min=1,max=255"`
	ModelName *string   `json:"model_name,omitempty"` // e.g., "gpt-4o-mini". Defaults to a system-wide default if not provided.
	UserID    uuid.UUID `json:"-"`                    // This will be set from the context, not the request body.
}

// PostMessageRequest represents the request to post a new message to a conversation.
type PostMessageRequest struct {
	Content   string                  `json:"content" validate:"required"`
	ModelID   *uuid.UUID              `json:"model_id,omitempty"` // Optional: use specific model for this message
	Artifacts []CreateArtifactRequest `json:"artifacts,omitempty"`
}

// CreateArtifactRequest represents the data needed to create a new artifact with a message.
type CreateArtifactRequest struct {
	Title    string `json:"title" validate:"required,min=1,max=255"`
	Type     string `json:"type" validate:"required"`
	Language *string `json:"language,omitempty"`
	Content  string `json:"content" validate:"required"`
}

// --- Response DTOs ---

// ConversationSummaryResponse represents a single conversation in a list.
type ConversationSummaryResponse struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	LastMessageAt *time.Time `json:"last_message_at"`
	ModelID       uuid.UUID  `json:"model_id"`
}

// MessageResponse represents a single message in a conversation.
type MessageResponse struct {
	ID        uuid.UUID       `json:"id"`
	Role      string          `json:"role"`
	Content   string          `json:"content"`
	CreatedAt time.Time       `json:"created_at"`
	Artifacts []ArtifactResponse `json:"artifacts,omitempty"`
}

// ArtifactResponse represents a single artifact in an API response.
type ArtifactResponse struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Type     string    `json:"type"`
	Language *string   `json:"language,omitempty"`
	Content  string    `json:"content"`
}

// ConversationDetailResponse represents the full details of a conversation.
type ConversationDetailResponse struct {
	ID           uuid.UUID         `json:"id"`
	Title        string            `json:"title"`
	ModelID      uuid.UUID         `json:"model_id"`
	SystemPrompt *string           `json:"system_prompt"`
	Messages     []MessageResponse `json:"messages"`
}

// JSONB is a local alias for map[string]interface{} for DTOs.
type JSONB map[string]interface{}

// ToolResponse represents a single tool in an API response.
type ToolResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Schema      JSONB     `json:"schema,omitempty"`
} 