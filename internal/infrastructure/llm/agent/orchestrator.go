package agent

import (
	"context"
	"fmt"
	"trading-alchemist/internal/domain/chat"
	"trading-alchemist/internal/domain/services"
	"trading-alchemist/internal/infrastructure/llm/provider"
)

// OrchestratorService is an implementation of services.LLMService that delegates
// to the appropriate provider client based on the provider's name.
type OrchestratorService struct {
}

// NewOrchestratorService creates a new LLM orchestrator.
func NewOrchestratorService() *OrchestratorService {
	return &OrchestratorService{}
}

// StreamChatCompletion finds the correct provider client and delegates the call.
func (s *OrchestratorService) StreamChatCompletion(
	ctx context.Context,
	providerE *chat.Provider,
	model *chat.Model,
	messages []*chat.Message,
	apiKey string,
	apiBaseOverride string,
) (<-chan services.ChatStreamEvent, error) {
	client, err := provider.NewClientForProvider(providerE.Name, apiKey, apiBaseOverride)
	if err != nil {
		return nil, fmt.Errorf("failed to create client for provider %s: %w", providerE.Name, err)
	}

	return client.StreamChatCompletion(ctx, model, messages)
} 