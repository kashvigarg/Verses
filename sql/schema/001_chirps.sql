-- goose Up
CREATE TABLE chirps(id INT,name VARCHAR(100) ,body TEXT,author_id UUID,created_at TIMESTAMP,upadated_at TIMESTAMP);

-- goose Down
DROP TABLE chirps;