-- +goose Up
CREATE TABLE users(
    Name VARCHAR(100) NOT NULL, 
    Email VARCHAR(100) NOT NULL,
    passwd bytea NOT NULL,
    id uuid PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
    );

-- +goose Down
DROP TABLE users;