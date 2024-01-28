-- +goose Up
CREATE TABLE chirps(
    id INT NOT NULL,
    body TEXT NOT NULL,
    author_id uuid REFERENCES users(id),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
    );

-- +goose Down
DROP TABLE chirps;

