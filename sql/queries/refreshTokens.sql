-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(
    created_at, 
    updated_at,
    user_id, 
    token,
    expires_at 
) VALUES (
    NOW(),
    NOW(),
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM refresh_tokens 
JOIN users ON users.id = refresh_tokens.user_id 
WHERE refresh_tokens.token = $1 
AND refresh_tokens.expires_at > NOW()
AND refresh_tokens.revoked_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = $1,
    updated_at = NOW()
WHERE token = $2;