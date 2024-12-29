-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, user_id) 
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;