-- name: UpsertUser :one
INSERT INTO users (github_id, username, avatar_url, email)
VALUES (?, ?, ?, ?)
ON CONFLICT (github_id)
DO UPDATE SET username = EXCLUDED.username, avatar_url = EXCLUDED.avatar_url, email = EXCLUDED.email, updated_at = datetime('now')
RETURNING *;

-- name: GetUserByEmail :one
SELECT id, username, email, unsubscribed, unsubscribe_token FROM users WHERE email = ?;

-- name: UnsubscribeUser :execrows
UPDATE users SET unsubscribed = true WHERE unsubscribe_token = ? AND unsubscribed = false;
