package entities

import (
	"time"

	"github.com/google/uuid"
)

// Provider represents LLM providers (OpenAI, Anthropic, Bedrock, etc.)
type Provider struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`                 // "openai", "anthropic", "bedrock"
	DisplayName string    `json:"display_name" db:"display_name"` // "OpenAI", "Anthropic Claude"
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Models      []*Model
} 