-- +goose Up
CREATE TABLE users (
    id                INTEGER PRIMARY KEY,
    github_id         INTEGER NOT NULL UNIQUE,
    username          TEXT NOT NULL,
    avatar_url        TEXT NOT NULL DEFAULT '',
    created_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email             TEXT NOT NULL DEFAULT '',
    unsubscribed      BOOLEAN NOT NULL DEFAULT false,
    unsubscribe_token TEXT NOT NULL DEFAULT (hex(randomblob(16)))
);

CREATE TABLE sessions (
    id         INTEGER PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT NOT NULL UNIQUE,
    csrf_token TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_token ON sessions (token);

CREATE TABLE comments (
    id           INTEGER PRIMARY KEY,
    article_slug TEXT NOT NULL,
    user_id      INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body         TEXT NOT NULL,
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    parent_id    INTEGER REFERENCES comments(id) ON DELETE CASCADE,
    approved     BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX idx_comments_article_slug ON comments (article_slug);

CREATE TABLE page_views (
    id           INTEGER PRIMARY KEY,
    article_slug TEXT NOT NULL,
    ip_address   TEXT NOT NULL DEFAULT '',
    user_agent   TEXT NOT NULL DEFAULT '',
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_page_views_article_slug ON page_views (article_slug);

-- +goose Down
DROP TABLE page_views;
DROP TABLE comments;
DROP TABLE sessions;
DROP TABLE users;
