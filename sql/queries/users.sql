-- name: Createuser :one
INSERT INTO users(Email, password, author_id, created_at) VALUES($1,$2,$3,$4)
RETURNING *;