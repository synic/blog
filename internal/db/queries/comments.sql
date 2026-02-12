-- name: CreateComment :one
INSERT INTO comments (article_slug, user_id, body, parent_id, approved)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListCommentsBySlug :many
SELECT c.id, c.article_slug, c.body, c.created_at, c.parent_id,
       u.username, u.avatar_url
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.article_slug = $1 AND c.approved = true
ORDER BY c.created_at ASC;

-- name: GetCommentWithUser :one
SELECT c.id, c.article_slug, c.body, c.created_at, c.parent_id, c.approved,
       u.username, u.avatar_url, u.email, u.unsubscribe_token, u.unsubscribed
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.id = $1;

-- name: CountCommentsBySlug :many
SELECT article_slug, count(*)::int AS comment_count
FROM comments
WHERE approved = true
GROUP BY article_slug;

-- name: ApproveComment :one
UPDATE comments SET approved = true WHERE id = $1 RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1;
