-- name: CreateTool :one
INSERT INTO tools (name, description, schema, provider_id, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, description, schema, provider_id, is_active, created_at, updated_at;

-- name: GetToolByID :one
SELECT id, name, description, schema, provider_id, is_active, created_at, updated_at FROM tools
WHERE id = $1;

-- name: GetToolByName :one
SELECT id, name, description, schema, provider_id, is_active, created_at, updated_at FROM tools
WHERE name = $1;

-- name: GetAvailableTools :many
SELECT id, name, description, schema, provider_id, is_active, created_at, updated_at FROM tools
WHERE is_active = true AND (provider_id IS NULL OR provider_id = $1)
ORDER BY name;

-- name: UpdateTool :one
UPDATE tools
SET
    description = $2,
    schema = $3,
    provider_id = $4,
    is_active = $5
WHERE id = $1
RETURNING id, name, description, schema, provider_id, is_active, created_at, updated_at;

-- name: DeleteTool :exec
DELETE FROM tools
WHERE id = $1;

-- name: LogToolUsage :one
INSERT INTO message_tools (message_id, tool_id, input, output, executed_at, duration, success, error)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, message_id, tool_id, input, output, executed_at, duration, success, error; 