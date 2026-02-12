-- +goose Up
CREATE TABLE comments (
    id           BIGSERIAL PRIMARY KEY,
    article_slug TEXT NOT NULL,
    user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body         TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_comments_article_slug ON comments (article_slug);

-- +goose Down
DROP TABLE comments;
