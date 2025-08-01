// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: user_provider_settings.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUserProviderSetting = `-- name: CreateUserProviderSetting :one
INSERT INTO user_provider_settings (user_id, provider_id, encrypted_api_key, api_base_override, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, provider_id, encrypted_api_key, api_base_override, is_active, created_at, updated_at
`

type CreateUserProviderSettingParams struct {
	UserID          pgtype.UUID `json:"user_id"`
	ProviderID      pgtype.UUID `json:"provider_id"`
	EncryptedApiKey pgtype.Text `json:"encrypted_api_key"`
	ApiBaseOverride pgtype.Text `json:"api_base_override"`
	IsActive        pgtype.Bool `json:"is_active"`
}

func (q *Queries) CreateUserProviderSetting(ctx context.Context, arg CreateUserProviderSettingParams) (UserProviderSetting, error) {
	row := q.db.QueryRow(ctx, createUserProviderSetting,
		arg.UserID,
		arg.ProviderID,
		arg.EncryptedApiKey,
		arg.ApiBaseOverride,
		arg.IsActive,
	)
	var i UserProviderSetting
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProviderID,
		&i.EncryptedApiKey,
		&i.ApiBaseOverride,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUserProviderSetting = `-- name: DeleteUserProviderSetting :exec
DELETE FROM user_provider_settings
WHERE id = $1
`

func (q *Queries) DeleteUserProviderSetting(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteUserProviderSetting, id)
	return err
}

const getUserProviderSetting = `-- name: GetUserProviderSetting :one
SELECT id, user_id, provider_id, encrypted_api_key, api_base_override, is_active, created_at, updated_at FROM user_provider_settings
WHERE user_id = $1 AND provider_id = $2
`

type GetUserProviderSettingParams struct {
	UserID     pgtype.UUID `json:"user_id"`
	ProviderID pgtype.UUID `json:"provider_id"`
}

func (q *Queries) GetUserProviderSetting(ctx context.Context, arg GetUserProviderSettingParams) (UserProviderSetting, error) {
	row := q.db.QueryRow(ctx, getUserProviderSetting, arg.UserID, arg.ProviderID)
	var i UserProviderSetting
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProviderID,
		&i.EncryptedApiKey,
		&i.ApiBaseOverride,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUserProviderSettings = `-- name: ListUserProviderSettings :many
SELECT id, user_id, provider_id, encrypted_api_key, api_base_override, is_active, created_at, updated_at FROM user_provider_settings
WHERE user_id = $1 AND is_active = true
ORDER BY created_at DESC
`

func (q *Queries) ListUserProviderSettings(ctx context.Context, userID pgtype.UUID) ([]UserProviderSetting, error) {
	rows, err := q.db.Query(ctx, listUserProviderSettings, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserProviderSetting{}
	for rows.Next() {
		var i UserProviderSetting
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ProviderID,
			&i.EncryptedApiKey,
			&i.ApiBaseOverride,
			&i.IsActive,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserProviderSetting = `-- name: UpdateUserProviderSetting :one
UPDATE user_provider_settings
SET
    encrypted_api_key = $2,
    api_base_override = $3,
    is_active = $4
WHERE id = $1
RETURNING id, user_id, provider_id, encrypted_api_key, api_base_override, is_active, created_at, updated_at
`

type UpdateUserProviderSettingParams struct {
	ID              pgtype.UUID `json:"id"`
	EncryptedApiKey pgtype.Text `json:"encrypted_api_key"`
	ApiBaseOverride pgtype.Text `json:"api_base_override"`
	IsActive        pgtype.Bool `json:"is_active"`
}

func (q *Queries) UpdateUserProviderSetting(ctx context.Context, arg UpdateUserProviderSettingParams) (UserProviderSetting, error) {
	row := q.db.QueryRow(ctx, updateUserProviderSetting,
		arg.ID,
		arg.EncryptedApiKey,
		arg.ApiBaseOverride,
		arg.IsActive,
	)
	var i UserProviderSetting
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProviderID,
		&i.EncryptedApiKey,
		&i.ApiBaseOverride,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
