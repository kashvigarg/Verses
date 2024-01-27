-- goose Up
CREATE TABLE users(name VARCHAR(100), Email VARCHAR(100),author_id UUID,passwd bytea,created_at TIMESTAMP,upadated_at TIMESTAMP);

-- goose Down
DROP TABLE users;