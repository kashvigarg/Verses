-- +goose Up
ALTER TABLE notifications ADD CONSTRAINT unique_comment_notification UNIQUE (user_id, prose_id, type, read);

-- +goose Down
ALTER TABLE notifications DROP CONSTRAINT unique_comment_notification;
