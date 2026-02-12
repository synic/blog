-- +goose Up
ALTER TABLE comments ADD COLUMN approved BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE users ADD COLUMN unsubscribed BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN unsubscribe_token TEXT NOT NULL DEFAULT gen_random_uuid()::text;

-- +goose Down
ALTER TABLE comments DROP COLUMN approved;
ALTER TABLE users DROP COLUMN unsubscribed;
ALTER TABLE users DROP COLUMN unsubscribe_token;
