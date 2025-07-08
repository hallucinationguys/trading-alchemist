package chat

import (
	"time"

	"github.com/google/uuid"
)

// UserProviderSetting stores user-specific settings for an LLM provider.
type UserProviderSetting struct {
	ID                uuid.UUID `json:"id" db:"id"`
	UserID            uuid.UUID `json:"user_id" db:"user_id"`
	ProviderID        uuid.UUID `json:"provider_id" db:"provider_id"`
	EncryptedAPIKey   *string   `json:"-" db:"encrypted_api_key"` // Not exposed in JSON responses
	APIBaseOverride   *string   `json:"api_base_override" db:"api_base_override"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
} 