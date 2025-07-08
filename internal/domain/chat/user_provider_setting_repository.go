package chat

import (
	"context"

	"github.com/google/uuid"
)

type UserProviderSettingRepository interface {
	Create(ctx context.Context, setting *UserProviderSetting) (*UserProviderSetting, error)
	GetByUserIDAndProviderID(ctx context.Context, userID, providerID uuid.UUID) (*UserProviderSetting, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*UserProviderSetting, error)
	Update(ctx context.Context, setting *UserProviderSetting) (*UserProviderSetting, error)
	Delete(ctx context.Context, id uuid.UUID) error
} 