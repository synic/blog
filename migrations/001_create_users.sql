-- +goose Up
CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    github_id  BIGINT NOT NULL UNIQUE,
    username   TEXT NOT NULL,
    avatar_url TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE users;
