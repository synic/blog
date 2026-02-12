-- +goose Up
ALTER TABLE comments ADD COLUMN parent_id BIGINT REFERENCES comments(id) ON DELETE CASCADE;
ALTER TABLE users ADD COLUMN email TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE comments DROP COLUMN parent_id;
ALTER TABLE users DROP COLUMN email;
