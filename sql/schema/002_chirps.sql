-- +goose Up
CREATE TABLE chirps(id INT NOT NULL, body TEXT,author_id UUID NOT NULL,created_at TIMESTAMP NOT NULL,upadated_at TIMESTAMP);

-- +goose Down
DROP TABLE chirps;