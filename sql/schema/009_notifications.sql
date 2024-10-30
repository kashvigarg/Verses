-- +goose Up
CREATE TABLE notifications(
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    actors VARCHAR [] ,
    generated_at TIMESTAMP NOT NULL,
    type VARCHAR(12) NOT NULL,
    read BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX sorted_notifications ON notifications(generated_at DESC);

-- +goose Down
DROP TABLE notifications;
DROP INDEX sorted_notifications;