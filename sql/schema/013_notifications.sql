-- +goose Up
ALTER TABLE notifications DROP CONSTRAINT unique_comment_notification;
ALTER TABLE notifications ADD CONSTRAINT unique_comment_notification UNIQUE (user_id,prose_id ,actors, type, read);

-- +goose Down
ALTER TABLE notifications ADD CONSTRAINT  unique_comment_notification UNIQUE (user_id, actors, type, read);
ALTER TABLE notifications DROP CONSTRAINT unique_comment_notification UNIQUE (user_id,prose_id ,actors, type, read);
