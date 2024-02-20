-- +goose Up
CREATE TABLE revocation(
    token bytea NOT NULL,
    revoked_at TIMESTAMP NOT NULL
    );

-- +goose Down
DROP TABLE revocation;