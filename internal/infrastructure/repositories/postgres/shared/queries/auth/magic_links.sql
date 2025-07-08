-- name: CreateMagicLink :one
INSERT INTO magic_links (user_id, token, token_hash, expires_at, ip_address, user_agent, purpose)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetMagicLinkByToken :one
SELECT ml.*, u.email, u.email_verified 
FROM magic_links ml
JOIN users u ON ml.user_id = u.id
WHERE ml.token = $1 
  AND ml.expires_at > NOW()
  AND ml.used_at IS NULL
  AND u.is_active = true;

-- name: UseMagicLink :one
UPDATE magic_links 
SET used_at = NOW()
WHERE id = $1
RETURNING *;

-- name: InvalidateUserMagicLinks :exec
UPDATE magic_links 
SET used_at = NOW()
WHERE user_id = $1 
  AND purpose = $2 
  AND used_at IS NULL;

-- name: CleanupExpiredMagicLinks :exec
DELETE FROM magic_links 
WHERE expires_at < NOW() - INTERVAL '1 day'; 