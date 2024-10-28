-- +goose Up
CREATE TABLE users(
    Name VARCHAR(100) NOT NULL, 
    Email VARCHAR(100) NOT NULL,
    passwd bytea NOT NULL,
    username VARCHAR(12) NOT NULL UNIQUE,
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
    );


-- +goose Down
DROP TABLE users;