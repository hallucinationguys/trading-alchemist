package chat

import (
	"time"

	"github.com/google/uuid"
)

// Model represents a language model available from a provider.
type Model struct {
	ID                uuid.UUID
	ProviderID        uuid.UUID
	Name              string
	DisplayName       string
	SupportsFunctions bool
	SupportsVision    bool
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
} 