-- name: CreateMessage :one
INSERT INTO messages (conversation_id, parent_id, role, content, model_id, token_count, cost, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, conversation_id, parent_id, role, content, model_id, token_count, cost, metadata, created_at, updated_at;

-- name: GetMessageByID :one
SELECT id, conversation_id, parent_id, role, content, model_id, token_count, cost, metadata, created_at, updated_at FROM messages
WHERE id = $1;

-- name: GetMessagesByConversationID :many
SELECT id, conversation_id, parent_id, role, content, model_id, token_count, cost, metadata, created_at, updated_at FROM messages
WHERE conversation_id = $1
ORDER BY created_at ASC
LIMIT $2 OFFSET $3;

-- name: CountMessagesByConversationID :one
SELECT COUNT(*) FROM messages
WHERE conversation_id = $1;

-- name: GetMessagesByConversationIDWithCursor :many
SELECT id, conversation_id, parent_id, role, content, model_id, token_count, cost, metadata, created_at, updated_at FROM messages
WHERE conversation_id = $1 AND created_at < $2
ORDER BY created_at DESC
LIMIT $3;

-- name: GetMessageThread :many
SELECT id, conversation_id, parent_id, role, content, model_id, token_count, cost, metadata, created_at, updated_at FROM messages
WHERE parent_id = $1
ORDER BY created_at ASC;

-- name: UpdateMessage :one
UPDATE messages
SET
    content = $2,
    token_count = $3,
    cost = $4,
    metadata = $5
WHERE id = $1
RETURNING id, conversation_id, parent_id, role, content, model_id, token_count, cost, metadata, created_at, updated_at;

-- name: DeleteMessage :exec
DELETE FROM messages
WHERE id = $1; 