-- +goose Up
ALTER TABLE users ADD Email varchar(90);

ALTER TABLE users ADD COLUMN created_at TIMESTAMP;

-- +goose Down
ALTER TABLE users DROP COLUMN Email;

ALTER TABLE users DROP COLUMN created_at;