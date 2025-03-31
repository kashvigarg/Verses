-- +goose Up
ALTER TABLE notifications 
ALTER COLUMN type TYPE VARCHAR(15);

-- goose Down
ALTER TABLE notifications 
ALTER COLUMN type TYPE VARCHAR(12);