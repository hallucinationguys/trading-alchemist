-- name: CreateProvider :one
INSERT INTO providers (name, display_name, is_active)
VALUES ($1, $2, $3)
RETURNING id, name, display_name, is_active, created_at, updated_at;

-- name: GetProviderByID :one
SELECT id, name, display_name, is_active, created_at, updated_at
FROM providers
WHERE id = $1
LIMIT 1;

-- name: GetProviderByName :one
SELECT id, name, display_name, is_active, created_at, updated_at
FROM providers
WHERE name = $1
LIMIT 1;

-- name: GetAllProviders :many
SELECT id, name, display_name, is_active, created_at, updated_at
FROM providers
ORDER BY display_name;

-- name: GetActiveProviders :many
SELECT id, name, display_name, is_active, created_at, updated_at
FROM providers
WHERE is_active = TRUE
ORDER BY display_name;

-- name: UpdateProvider :one
UPDATE providers
SET
    display_name = $2,
    is_active = $3
WHERE id = $1
RETURNING id, name, display_name, is_active, created_at, updated_at;

-- name: DeleteProvider :exec
DELETE FROM providers
WHERE id = $1;

-- name: GetProvidersWithModels :many
SELECT 
    p.id as provider_id,
    p.name as provider_name,
    p.display_name as provider_display_name,
    p.is_active as provider_is_active,
    p.created_at as provider_created_at,
    p.updated_at as provider_updated_at,
    m.id as model_id,
    m.name as model_name,
    m.display_name as model_display_name,
    m.supports_functions as model_supports_functions,
    m.supports_vision as model_supports_vision,
    m.is_active as model_is_active,
    m.created_at as model_created_at,
    m.updated_at as model_updated_at
FROM providers p
LEFT JOIN models m ON p.id = m.provider_id AND m.is_active = TRUE
WHERE p.is_active = TRUE
ORDER BY p.display_name, m.display_name;

-- name: GetAvailableModelsForUser :many
SELECT 
    p.id as provider_id,
    p.name as provider_name,
    p.display_name as provider_display_name,
    m.id as model_id,
    m.name as model_name,
    m.display_name as model_display_name,
    m.supports_functions as model_supports_functions,
    m.supports_vision as model_supports_vision,
    ups.id as setting_id,
    ups.encrypted_api_key as has_api_key,
    ups.is_active as setting_is_active
FROM providers p
INNER JOIN models m ON p.id = m.provider_id AND m.is_active = TRUE
LEFT JOIN user_provider_settings ups ON p.id = ups.provider_id AND ups.user_id = $1
WHERE p.is_active = TRUE
ORDER BY 
    CASE WHEN ups.encrypted_api_key IS NOT NULL AND ups.is_active = TRUE THEN 0 ELSE 1 END,
    p.display_name,
    m.display_name; 