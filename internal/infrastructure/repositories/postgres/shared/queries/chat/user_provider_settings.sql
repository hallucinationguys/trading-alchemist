-- name: CreateUserProviderSetting :one
INSERT INTO user_provider_settings (user_id, provider_id, encrypted_api_key, api_base_override, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, provider_id, encrypted_api_key, api_base_override, is_active, created_at, updated_at;

-- name: GetUserProviderSetting :one
SELECT id, user_id, provider_id, encrypted_api_key, api_base_override, is_active, created_at, updated_at FROM user_provider_settings
WHERE user_id = $1 AND provider_id = $2;

-- name: ListUserProviderSettings :many
SELECT id, user_id, provider_id, encrypted_api_key, api_base_override, is_active, created_at, updated_at FROM user_provider_settings
WHERE user_id = $1 AND is_active = true
ORDER BY created_at DESC;

-- name: UpdateUserProviderSetting :one
UPDATE user_provider_settings
SET
    encrypted_api_key = $2,
    api_base_override = $3,
    is_active = $4
WHERE id = $1
RETURNING id, user_id, provider_id, encrypted_api_key, api_base_override, is_active, created_at, updated_at;

-- name: DeleteUserProviderSetting :exec
DELETE FROM user_provider_settings
WHERE id = $1; 