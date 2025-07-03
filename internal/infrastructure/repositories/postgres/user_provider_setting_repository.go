package postgres

import (
	"context"
	"fmt"
	"trading-alchemist/internal/domain/entities"
	"trading-alchemist/internal/domain/repositories"
	"trading-alchemist/internal/infrastructure/repositories/postgres/sqlc"
	"trading-alchemist/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserProviderSettingRepository struct {
	queries *sqlc.Queries
}

func NewUserProviderSettingRepository(db sqlc.DBTX) repositories.UserProviderSettingRepository {
	return &UserProviderSettingRepository{
		queries: sqlc.New(db),
	}
}

func (r *UserProviderSettingRepository) Create(ctx context.Context, setting *entities.UserProviderSetting) (*entities.UserProviderSetting, error) {
	params := sqlc.CreateUserProviderSettingParams{
		UserID:     pgtype.UUID{Bytes: setting.UserID, Valid: true},
		ProviderID: pgtype.UUID{Bytes: setting.ProviderID, Valid: true},
		IsActive:   pgtype.Bool{Bool: setting.IsActive, Valid: true},
	}
	if setting.EncryptedAPIKey != nil {
		params.EncryptedApiKey = pgtype.Text{String: *setting.EncryptedAPIKey, Valid: true}
	}
	if setting.APIBaseOverride != nil {
		params.ApiBaseOverride = pgtype.Text{String: *setting.APIBaseOverride, Valid: true}
	}

	sqlcSetting, err := r.queries.CreateUserProviderSetting(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user provider setting: %w", err)
	}
	return sqlcUserProviderSettingToEntity(&sqlcSetting), nil
}

func (r *UserProviderSettingRepository) GetByUserIDAndProviderID(ctx context.Context, userID, providerID uuid.UUID) (*entities.UserProviderSetting, error) {
	params := sqlc.GetUserProviderSettingParams{
		UserID:     pgtype.UUID{Bytes: userID, Valid: true},
		ProviderID: pgtype.UUID{Bytes: providerID, Valid: true},
	}
	sqlcSetting, err := r.queries.GetUserProviderSetting(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrUserProviderSettingNotFound
		}
		return nil, fmt.Errorf("failed to get user provider setting: %w", err)
	}
	return sqlcUserProviderSettingToEntity(&sqlcSetting), nil
}

func (r *UserProviderSettingRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.UserProviderSetting, error) {
	userUUID := pgtype.UUID{Bytes: userID, Valid: true}
	sqlcSettings, err := r.queries.ListUserProviderSettings(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user provider settings: %w", err)
	}
	settings := make([]*entities.UserProviderSetting, len(sqlcSettings))
	for i, s := range sqlcSettings {
		settings[i] = sqlcUserProviderSettingToEntity(&s)
	}
	return settings, nil
}

func (r *UserProviderSettingRepository) Update(ctx context.Context, setting *entities.UserProviderSetting) (*entities.UserProviderSetting, error) {
	params := sqlc.UpdateUserProviderSettingParams{
		ID:       pgtype.UUID{Bytes: setting.ID, Valid: true},
		IsActive: pgtype.Bool{Bool: setting.IsActive, Valid: true},
	}
	if setting.EncryptedAPIKey != nil {
		params.EncryptedApiKey = pgtype.Text{String: *setting.EncryptedAPIKey, Valid: true}
	}
	if setting.APIBaseOverride != nil {
		params.ApiBaseOverride = pgtype.Text{String: *setting.APIBaseOverride, Valid: true}
	}

	sqlcSetting, err := r.queries.UpdateUserProviderSetting(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrUserProviderSettingNotFound
		}
		return nil, fmt.Errorf("failed to update user provider setting: %w", err)
	}
	return sqlcUserProviderSettingToEntity(&sqlcSetting), nil
}

func (r *UserProviderSettingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	settingUUID := pgtype.UUID{Bytes: id, Valid: true}
	return r.queries.DeleteUserProviderSetting(ctx, settingUUID)
}

func sqlcUserProviderSettingToEntity(s *sqlc.UserProviderSetting) *entities.UserProviderSetting {
	setting := &entities.UserProviderSetting{
		IsActive:  s.IsActive.Bool,
		CreatedAt: s.CreatedAt.Time,
		UpdatedAt: s.UpdatedAt.Time,
	}
	if s.ID.Valid {
		setting.ID = s.ID.Bytes
	}
	if s.UserID.Valid {
		setting.UserID = s.UserID.Bytes
	}
	if s.ProviderID.Valid {
		setting.ProviderID = s.ProviderID.Bytes
	}
	if s.EncryptedApiKey.Valid {
		setting.EncryptedAPIKey = &s.EncryptedApiKey.String
	}
	if s.ApiBaseOverride.Valid {
		setting.APIBaseOverride = &s.ApiBaseOverride.String
	}
	return setting
} 