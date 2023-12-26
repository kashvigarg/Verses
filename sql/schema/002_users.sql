-- +goose Up
CREATE TABLE users(id SMALLINT,body TEXT,Author_id UUID);

-- +goose Down
DROP TABLE users;

