-- name: CreateSession :one
INSERT INTO sessions (user_id, token, csrf_token, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSessionByToken :one
SELECT s.id, s.user_id, s.token, s.csrf_token, s.expires_at, u.username, u.avatar_url, u.email
FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.token = $1 AND s.expires_at > now();

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token = $1;

-- name: DeleteExpiredSessions :execrows
DELETE FROM sessions WHERE expires_at < now();
