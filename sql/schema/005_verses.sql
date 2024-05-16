-- +goose Up
ALTER TABLE chirps RENAME TO prose;

-- +goose Down
ALTER TABLE prose RENAME TO chirps;