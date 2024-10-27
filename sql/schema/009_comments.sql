-- +goose Up
CREATE TABLE comments(
    id SERIAL PRIMARY KEY,
    prose_id uuid REFERENCES prose(id),
    user_id uuid REFERENCES users(id),
    body TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    likes_count INT NOT NULL DEFAULT 0 CHECK (likes_count>=0)
);

CREATE INDEX post_comments ON comments(created_at DESC);

CREATE TABLE comment_likes(
    comment_id INT NOT NULL REFERENCES comments(id),
    user_id uuid NOT NULL REFERENCES users(id),
    PRIMARY KEY (prose_id,user_id)
);

ALTER TABLE prose
ADD COLUMN comments INT NOT NULL DEFAULT 0 CHECK (comments>=0);

-- +goose Down
DROP TABLE comments;

DROP INDEX post_comments;

DROP TABLE comment_likes;

ALTER TABLE prose
DROP COLUMN comments;