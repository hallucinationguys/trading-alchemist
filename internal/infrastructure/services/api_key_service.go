package services

import (
	"context"
	"fmt"
	"trading-alchemist/internal/domain/chat"
	"trading-alchemist/internal/domain/services"
	"trading-alchemist/pkg/errors"

	"github.com/google/uuid"
)

// APIKeyServiceImpl implements the APIKeyService interface
type APIKeyServiceImpl struct {
	userProviderRepo chat.UserProviderSettingRepository
}

// NewAPIKeyService creates a new API key service
func NewAPIKeyService(userProviderRepo chat.UserProviderSettingRepository) services.APIKeyService {
	return &APIKeyServiceImpl{
		userProviderRepo: userProviderRepo,
	}
}

// IsProviderConfigured checks if a user has an active API key for a provider
func (s *APIKeyServiceImpl) IsProviderConfigured(ctx context.Context, userID, providerID uuid.UUID) (bool, error) {
	setting, err := s.userProviderRepo.GetByUserIDAndProviderID(ctx, userID, providerID)
	if err != nil {
		if err == errors.ErrUserProviderSettingNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to check provider configuration: %w", err)
	}
	
	// Provider is configured if it has an API key and is active
	isConfigured := setting.EncryptedAPIKey != nil && 
		*setting.EncryptedAPIKey != "" && 
		setting.IsActive
	
	return isConfigured, nil
}

// GetUserProviderConfig returns the configuration for a user-provider combination
func (s *APIKeyServiceImpl) GetUserProviderConfig(ctx context.Context, userID, providerID uuid.UUID) (*chat.UserProviderSetting, error) {
	setting, err := s.userProviderRepo.GetByUserIDAndProviderID(ctx, userID, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider configuration: %w", err)
	}
	return setting, nil
}

// ValidateProviderAccess checks if user can access a specific provider
func (s *APIKeyServiceImpl) ValidateProviderAccess(ctx context.Context, userID, providerID uuid.UUID) error {
	isConfigured, err := s.IsProviderConfigured(ctx, userID, providerID)
	if err != nil {
		return fmt.Errorf("failed to validate provider access: %w", err)
	}
	
	if !isConfigured {
		return errors.NewAppError(
			errors.CodeConfiguration, 
			"API key for this provider is not configured or not active", 
			nil,
		)
	}
	
	return nil
}

// GetConfiguredProviders returns all providers that the user has configured API keys for
func (s *APIKeyServiceImpl) GetConfiguredProviders(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	settings, err := s.userProviderRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user provider settings: %w", err)
	}
	
	var configuredProviders []uuid.UUID
	for _, setting := range settings {
		if setting.EncryptedAPIKey != nil && 
			*setting.EncryptedAPIKey != "" && 
			setting.IsActive {
			configuredProviders = append(configuredProviders, setting.ProviderID)
		}
	}
	
	return configuredProviders, nil
} 