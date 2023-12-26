-- +goose Up
ALTER TABLE users DROP body;

ALTER TABLE users ADD COLUMN password bytea;

-- +goose Down
ALTER TABLE users ADD COLUMN body TEXT;

ALTER TABLE users DROP COLUMN password;