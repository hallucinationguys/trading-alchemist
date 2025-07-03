package repositories

import (
	"context"
	"trading-alchemist/internal/domain/entities"

	"github.com/google/uuid"
)

type UserProviderSettingRepository interface {
	Create(ctx context.Context, setting *entities.UserProviderSetting) (*entities.UserProviderSetting, error)
	GetByUserIDAndProviderID(ctx context.Context, userID, providerID uuid.UUID) (*entities.UserProviderSetting, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.UserProviderSetting, error)
	Update(ctx context.Context, setting *entities.UserProviderSetting) (*entities.UserProviderSetting, error)
	Delete(ctx context.Context, id uuid.UUID) error
} 