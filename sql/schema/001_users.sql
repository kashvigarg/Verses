-- +goose Up
CREATE TABLE users(Email varchar(100) NOT NULL, passwd bytea NOT NULL,id UUID NOT NULL,red BOOL DEFAULT NULL,created_at TIMESTAMP NOT NULL,updated_at TIMESTAMP NOT NULL);

-- +goose Down
DROP TABLE users;