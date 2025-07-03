package agent

import (
	"trading-alchemist/internal/domain/services"
)

// NewLLMService creates the central LLM orchestrator.
func NewLLMService() (services.LLMService, error) {
	orchestrator := NewOrchestratorService()
	return orchestrator, nil
} 