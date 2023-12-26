-- +goose Up
ALTER TABLE users DROP COLUMN id;

ALTER TABLE users ADD COLUMN chirpy_red BOOL;

-- +goose Down

ALTER TABLE users ADD COLUMN id SMALLINT;

ALTER TABLE users DROP COLUMN chirpy_red;