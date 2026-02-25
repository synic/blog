-- name: CreateComment :one
INSERT INTO comments (article_slug, user_id, body, parent_id, approved)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: ListCommentsBySlug :many
SELECT c.id, c.article_slug, c.body, c.created_at, c.parent_id,
       u.username, u.avatar_url
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.article_slug = ? AND c.approved = true
ORDER BY c.created_at ASC;

-- name: GetCommentWithUser :one
SELECT c.id, c.article_slug, c.body, c.created_at, c.parent_id, c.approved,
       u.username, u.avatar_url, u.email, u.unsubscribe_token, u.unsubscribed
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.id = ?;

-- name: CountCommentsBySlug :many
SELECT article_slug, CAST(count(*) AS int) AS comment_count
FROM comments
WHERE approved = true
GROUP BY article_slug;

-- name: ApproveComment :one
UPDATE comments SET approved = true WHERE id = ? RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = ?;
