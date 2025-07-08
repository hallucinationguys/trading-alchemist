-- name: CreateConversation :one
INSERT INTO conversations (user_id, title, model_id, system_prompt, settings)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, title, model_id, system_prompt, settings, is_archived, created_at, updated_at, last_message_at;

-- name: GetConversationByID :one
SELECT id, user_id, title, model_id, system_prompt, settings, is_archived, created_at, updated_at, last_message_at FROM conversations
WHERE id = $1;

-- name: GetConversationsByUserID :many
SELECT id, user_id, title, model_id, system_prompt, settings, is_archived, created_at, updated_at, last_message_at FROM conversations
WHERE user_id = $1 AND is_archived = false
ORDER BY last_message_at DESC NULLS LAST, created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateConversation :one
UPDATE conversations
SET
    title = $2,
    model_id = $3,
    system_prompt = $4,
    settings = $5
WHERE id = $1
RETURNING id, user_id, title, model_id, system_prompt, settings, is_archived, created_at, updated_at, last_message_at;

-- name: UpdateConversationLastMessageAt :exec
UPDATE conversations
SET last_message_at = $2
WHERE id = $1;

-- name: UpdateConversationTitle :exec
UPDATE conversations
SET title = $2
WHERE id = $1;

-- name: ArchiveConversation :exec
UPDATE conversations
SET is_archived = true
WHERE id = $1;

-- name: DeleteConversation :exec
DELETE FROM conversations
WHERE id = $1; 