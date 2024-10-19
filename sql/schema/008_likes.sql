-- +goose Up
CREATE TABLE post_likes(
    prose_id uuid NOT NULL REFERENCES prose(id),
    user_id uuid NOT NULL REFERENCES users(id),
    PRIMARY KEY (post_id,user_id)
);

ALTER TABLE prose
ADD COLUMN likes INT NOT NULL DEFAULT 0 CHECK (likes>=0);

-- +goose Down
DROP TABLE post_likes;

ALTER TABLE prose
DROP COLUMN likes;
