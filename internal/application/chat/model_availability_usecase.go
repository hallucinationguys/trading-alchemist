package chat

import (
	"context"
	"fmt"
	"trading-alchemist/internal/domain/chat"
	"trading-alchemist/internal/domain/services"
	"trading-alchemist/internal/infrastructure/database"

	"github.com/google/uuid"
)

// ModelAvailabilityUseCase handles fetching available models with API key status efficiently
type ModelAvailabilityUseCase struct {
	dbService     *database.Service
	apiKeyService services.APIKeyService
}

// NewModelAvailabilityUseCase creates a new ModelAvailabilityUseCase
func NewModelAvailabilityUseCase(
	dbService *database.Service,
	apiKeyService services.APIKeyService,
) *ModelAvailabilityUseCase {
	return &ModelAvailabilityUseCase{
		dbService:     dbService,
		apiKeyService: apiKeyService,
	}
}

// GetAvailableModelsWithAPIKeyStatus returns all available models with their API key configuration status
// This method uses optimized JOIN queries to eliminate N+1 database calls
func (uc *ModelAvailabilityUseCase) GetAvailableModelsWithAPIKeyStatus(ctx context.Context, userID uuid.UUID) ([]ProviderResponse, error) {
	var providers []*chat.Provider
	
	err := uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		var err error
		providers, err = provider.Provider().GetAvailableModelsForUser(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get available models for user: %w", err)
		}
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// Transform to DTO response
	result := make([]ProviderResponse, len(providers))
	for i, provider := range providers {
		models := make([]ModelResponse, len(provider.Models))
		for j, model := range provider.Models {
			// Check if user has API key for this provider
			hasAPIKey, _ := uc.apiKeyService.IsProviderConfigured(ctx, userID, provider.ID)
			
			// Create tags based on model capabilities and API key status
			tags := []string{"LLM", "CHAT"}
			if model.SupportsVision {
				tags = append(tags, "VISION")
			}
			if hasAPIKey {
				tags = append(tags, "CONFIGURED")
			} else {
				tags = append(tags, "NEEDS_API_KEY")
			}
			
			models[j] = ModelResponse{
				ID:          model.ID.String(),
				Name:        model.Name,
				DisplayName: model.DisplayName,
				IsActive:    model.IsActive,
				Tags:        tags,
			}
		}
		
		result[i] = ProviderResponse{
			ID:          provider.ID,
			Name:        provider.Name,
			DisplayName: provider.DisplayName,
			IsActive:    provider.IsActive,
			Models:      models,
		}
	}
	
	return result, nil
} 