-- name: UpsertUser :one
INSERT INTO users (github_id, username, avatar_url, email)
VALUES ($1, $2, $3, $4)
ON CONFLICT (github_id)
DO UPDATE SET username = EXCLUDED.username, avatar_url = EXCLUDED.avatar_url, email = EXCLUDED.email, updated_at = now()
RETURNING *;

-- name: GetUserByEmail :one
SELECT id, username, email, unsubscribed, unsubscribe_token FROM users WHERE email = $1;

-- name: UnsubscribeUser :execrows
UPDATE users SET unsubscribed = true WHERE unsubscribe_token = $1 AND unsubscribed = false;
