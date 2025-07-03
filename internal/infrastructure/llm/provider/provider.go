package provider

import (
	"context"
	"trading-alchemist/internal/domain/entities"
	"trading-alchemist/internal/domain/services"
)

// ProviderClient defines the interface that each LLM provider's client must implement.
type ProviderClient interface {
	StreamChatCompletion(
		ctx context.Context,
		model *entities.Model,
		messages []*entities.Message,
	) (<-chan services.ChatStreamEvent, error)
	// In the future, we could add other methods like:
	// GetToolDefinitions() []ToolDefinition
} 