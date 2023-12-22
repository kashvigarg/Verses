-- +goose Up
CREATE TABLE chirps(id SMALLINT,body TEXT,Author_id UUID);

-- +goose Down
DROP TABLE chirps;