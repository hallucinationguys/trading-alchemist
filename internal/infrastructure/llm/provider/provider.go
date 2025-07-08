package provider

import (
	"context"
	"trading-alchemist/internal/domain/chat"
	"trading-alchemist/internal/domain/services"
)

// ProviderClient defines the interface that each LLM provider's client must implement.
type ProviderClient interface {
	StreamChatCompletion(
		ctx context.Context,
		model *chat.Model,
		messages []*chat.Message,
	) (<-chan services.ChatStreamEvent, error)
	// In the future, we could add other methods like:
	// GetToolDefinitions() []ToolDefinition
} 