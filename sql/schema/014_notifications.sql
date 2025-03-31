-- +goose Up
ALTER TABLE notifications DROP CONSTRAINT notifications_prose_id_key;

-- +goose Down
ALTER TABLE notifications ADD CONSTRAINT notifications_prose_id_key UNIQUE (prose_id);
