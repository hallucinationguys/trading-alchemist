-- name: CreateUser :one
INSERT INTO users (email, first_name, last_name, email_verified)
VALUES ($1, $2, $3, $4)
RETURNING id, email, email_verified, first_name, last_name, avatar_url, is_active, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, email, email_verified, first_name, last_name, avatar_url, is_active, created_at, updated_at FROM users 
WHERE id = $1 AND is_active = true;

-- name: GetUserByEmail :one
SELECT id, email, email_verified, first_name, last_name, avatar_url, is_active, created_at, updated_at FROM users 
WHERE email = $1 AND is_active = true;

-- name: UpdateUser :one
UPDATE users 
SET first_name = $2, last_name = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, email, email_verified, first_name, last_name, avatar_url, is_active, created_at, updated_at;

-- name: VerifyUserEmail :one
UPDATE users 
SET email_verified = true, updated_at = NOW()
WHERE id = $1
RETURNING id, email, email_verified, first_name, last_name, avatar_url, is_active, created_at, updated_at;

-- name: DeactivateUser :exec
UPDATE users 
SET is_active = false, updated_at = NOW()
WHERE id = $1; 