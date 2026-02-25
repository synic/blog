-- name: CreatePageView :exec
INSERT INTO page_views (article_slug, ip_address, user_agent)
VALUES (?, ?, ?);

-- name: CountPageViewsBySlug :many
SELECT article_slug, CAST(count(*) AS int) AS view_count
FROM page_views
GROUP BY article_slug
ORDER BY view_count DESC;
