-- +goose Up
ALTER TABLE users 
ADD COLUMN followers INT NOT NULL DEFAULT 0 CHECK (followers>=0),
ADD COLUMN followees INT NOT NULL DEFAULT 0 CHECK (followees>=0);

CREATE TABLE follows(
    follower_id uuid NOT NULL REFERENCES users(id),
    followee_id uuid NOT NULL REFERENCES users(id),
    PRIMARY KEY (follower_id,followee_id) 
);

-- +goose Down
ALTER TABLE users 
DROP COLUMN followers, 
DROP COLUMN followees;

DROP TABLE follows;