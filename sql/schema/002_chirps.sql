-- +goose Up
CREATE TABLE chirps(
    id uuid PRIMARY KEY,
    body TEXT NOT NULL,
    author_id uuid NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
    );

-- +goose Down
DROP TABLE chirps;
