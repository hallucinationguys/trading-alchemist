package usecases

import (
	"context"
	"fmt"
	"trading-alchemist/internal/application/dto"
	"trading-alchemist/internal/config"
	"trading-alchemist/internal/domain/entities"
	"trading-alchemist/internal/infrastructure/database"
	"trading-alchemist/pkg/errors"
	"trading-alchemist/pkg/utils"

	"github.com/google/uuid"
)

// UserProviderSettingUseCase handles business logic for user provider settings.
type UserProviderSettingUseCase struct {
	dbService *database.Service
	config    *config.Config
}

// NewUserProviderSettingUseCase creates a new UserProviderSettingUseCase.
func NewUserProviderSettingUseCase(dbService *database.Service, config *config.Config) *UserProviderSettingUseCase {
	return &UserProviderSettingUseCase{
		dbService: dbService,
		config:    config,
	}
}

// ListProviders returns all available providers in the system.
func (uc *UserProviderSettingUseCase) ListProviders(ctx context.Context) ([]dto.ProviderResponse, error) {
	var providers []*entities.Provider
	err := uc.dbService.ExecuteInTx(ctx, func(providerRepo database.RepositoryProvider) error {
		var err error
		providers, err = providerRepo.Provider().GetActive(ctx)
		if err != nil {
			return err
		}

		for _, p := range providers {
			models, err := providerRepo.Model().GetActiveModelsByProviderID(ctx, p.ID)
			if err != nil {
				// Log the error but don't fail the entire operation
				fmt.Printf("Warning: failed to get models for provider %s: %v\n", p.Name, err)
				p.Models = []*entities.Model{}
			} else {
				p.Models = models
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list active providers: %w", err)
	}

	response := make([]dto.ProviderResponse, len(providers))
	for i, p := range providers {
		response[i] = dto.ToProviderResponse(p)
	}
	return response, nil
}

// ListUserSettings returns all of a user's configured provider settings.
func (uc *UserProviderSettingUseCase) ListUserSettings(ctx context.Context, userID uuid.UUID) ([]dto.UserProviderSettingResponse, error) {
	var settings []*entities.UserProviderSetting
	err := uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		var err error
		settings, err = provider.UserProviderSetting().ListByUserID(ctx, userID)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list user provider settings: %w", err)
	}

	response := make([]dto.UserProviderSettingResponse, len(settings))
	for i, s := range settings {
		// This requires another DB call to get provider details.
		// For performance, a JOIN in the query would be better, but for now, this is simpler.
		var p *entities.Provider
		err = uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
			var errTx error
			p, errTx = provider.Provider().GetByID(ctx, s.ProviderID)
			return errTx
		})
		if err != nil {
			// If a provider for a setting is not found, we can either skip it or return an error.
			// Skipping is safer for the user experience.
			continue
		}
		response[i] = dto.ToUserProviderSettingResponse(s, p)
	}

	return response, nil
}

// UpsertUserProviderSetting creates or updates a user's provider setting.
func (uc *UserProviderSettingUseCase) UpsertUserProviderSetting(ctx context.Context, userID uuid.UUID, req *dto.UpsertUserProviderSettingRequest) (*dto.UserProviderSettingResponse, error) {
	var setting *entities.UserProviderSetting
	var providerInfo *entities.Provider

	encryptionKey, err := uc.config.GetEncryptionKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get encryption key: %w", err)
	}

	err = uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		// Check if provider exists
		var errTx error
		providerInfo, errTx = provider.Provider().GetByID(ctx, req.ProviderID)
		if errTx != nil {
			if errTx == errors.ErrProviderNotFound {
				return errors.NewAppError(errors.CodeNotFound, "The specified provider does not exist.", errTx)
			}
			return fmt.Errorf("failed to verify provider: %w", errTx)
		}

		// Try to get existing setting
		existingSetting, errTx := provider.UserProviderSetting().GetByUserIDAndProviderID(ctx, userID, req.ProviderID)
		if errTx != nil && errTx != errors.ErrUserProviderSettingNotFound {
			return fmt.Errorf("failed to check for existing setting: %w", errTx)
		}

		if existingSetting != nil {
			// Update existing setting
			if req.APIKey != "" {
				// Only update API key if a new one is provided
				encryptedAPIKey, err := utils.Encrypt(req.APIKey, encryptionKey)
				if err != nil {
					return fmt.Errorf("failed to encrypt API key: %w", err)
				}
				existingSetting.EncryptedAPIKey = &encryptedAPIKey
			}
			// Update other fields
			existingSetting.APIBaseOverride = req.APIBaseOverride
			if req.IsActive != nil {
				existingSetting.IsActive = *req.IsActive
			}
			setting, errTx = provider.UserProviderSetting().Update(ctx, existingSetting)
		} else {
			// Create new setting - API key is required
			if req.APIKey == "" {
				return errors.NewAppError(errors.CodeValidation, "API key is required for new provider settings", nil)
			}
			
			encryptedAPIKey, err := utils.Encrypt(req.APIKey, encryptionKey)
			if err != nil {
				return fmt.Errorf("failed to encrypt API key: %w", err)
			}
			
			newSetting := &entities.UserProviderSetting{
				UserID:          userID,
				ProviderID:      req.ProviderID,
				EncryptedAPIKey: &encryptedAPIKey,
				APIBaseOverride: req.APIBaseOverride,
				IsActive:        true, // Default to active for new settings
			}
			// Only override if explicitly set in request
			if req.IsActive != nil {
				newSetting.IsActive = *req.IsActive
			}
			setting, errTx = provider.UserProviderSetting().Create(ctx, newSetting)
		}
		return errTx
	})

	if err != nil {
		return nil, err
	}

	response := dto.ToUserProviderSettingResponse(setting, providerInfo)
	return &response, nil
} 