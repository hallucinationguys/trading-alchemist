package chat

import (
	"time"
	"trading-alchemist/internal/domain/shared"

	"github.com/google/uuid"
)

// Tool represents MCP tools or function calls
type Tool struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Schema      shared.JSONB      `json:"schema" db:"schema"` // JSON schema for parameters
	ProviderID  *uuid.UUID `json:"provider_id" db:"provider_id"` // Optional: tool specific to provider
	IsActive    bool       `json:"is_active" db:"is_active"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
} 