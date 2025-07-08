-- name: CreateModel :one
INSERT INTO models (
    provider_id, name, display_name, supports_functions, supports_vision, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, provider_id, name, display_name, supports_functions, supports_vision, is_active, created_at, updated_at;

-- name: GetModelByID :one
SELECT id, provider_id, name, display_name, supports_functions, supports_vision, is_active, created_at, updated_at FROM models
WHERE id = $1
LIMIT 1;

-- name: GetModelByName :one
SELECT id, provider_id, name, display_name, supports_functions, supports_vision, is_active, created_at, updated_at FROM models
WHERE provider_id = $1 AND name = $2
LIMIT 1;

-- name: GetModelsByProviderID :many
SELECT id, provider_id, name, display_name, supports_functions, supports_vision, is_active, created_at, updated_at FROM models
WHERE provider_id = $1
ORDER BY display_name;

-- name: GetActiveModelsByProviderID :many
SELECT id, provider_id, name, display_name, supports_functions, supports_vision, is_active, created_at, updated_at FROM models
WHERE provider_id = $1 AND is_active = TRUE
ORDER BY name;

-- name: UpdateModel :one
UPDATE models
SET
    display_name = $2,
    supports_functions = $3,
    supports_vision = $4,
    is_active = $5
WHERE id = $1
RETURNING id, provider_id, name, display_name, supports_functions, supports_vision, is_active, created_at, updated_at;

-- name: DeleteModel :exec
DELETE FROM models
WHERE id = $1; 