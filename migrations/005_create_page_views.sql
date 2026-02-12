-- +goose Up
CREATE TABLE page_views (
    id BIGSERIAL PRIMARY KEY,
    article_slug TEXT NOT NULL,
    ip_address TEXT NOT NULL DEFAULT '',
    user_agent TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_page_views_article_slug ON page_views (article_slug);

-- +goose Down
DROP TABLE page_views;
