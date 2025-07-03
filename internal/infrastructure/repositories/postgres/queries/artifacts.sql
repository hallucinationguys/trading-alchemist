-- name: CreateArtifact :one
INSERT INTO artifacts (message_id, title, type, language, content, content_hash, size, is_public)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, message_id, title, type, language, content, content_hash, size, is_public, created_at, updated_at;

-- name: GetArtifactByID :one
SELECT id, message_id, title, type, language, content, content_hash, size, is_public, created_at, updated_at FROM artifacts
WHERE id = $1;

-- name: GetArtifactsByMessageID :many
SELECT id, message_id, title, type, language, content, content_hash, size, is_public, created_at, updated_at FROM artifacts
WHERE message_id = $1
ORDER BY created_at ASC;

-- name: GetPublicArtifacts :many
SELECT id, message_id, title, type, language, content, content_hash, size, is_public, created_at, updated_at FROM artifacts
WHERE is_public = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateArtifact :one
UPDATE artifacts
SET
    title = $2,
    content = $3,
    content_hash = $4,
    size = $5,
    is_public = $6
WHERE id = $1
RETURNING id, message_id, title, type, language, content, content_hash, size, is_public, created_at, updated_at;

-- name: DeleteArtifact :exec
DELETE FROM artifacts
WHERE id = $1; 