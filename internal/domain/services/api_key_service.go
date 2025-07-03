package services

import (
	"context"
	"trading-alchemist/internal/domain/entities"

	"github.com/google/uuid"
)

// APIKeyService encapsulates all API key-related business logic
type APIKeyService interface {
	// IsProviderConfigured checks if a user has an active API key for a provider
	IsProviderConfigured(ctx context.Context, userID, providerID uuid.UUID) (bool, error)
	
	// GetUserProviderConfig returns the configuration for a user-provider combination
	GetUserProviderConfig(ctx context.Context, userID, providerID uuid.UUID) (*entities.UserProviderSetting, error)
	
	// ValidateProviderAccess checks if user can access a specific provider
	ValidateProviderAccess(ctx context.Context, userID, providerID uuid.UUID) error
	
	// GetConfiguredProviders returns all providers that the user has configured API keys for
	GetConfiguredProviders(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
} 