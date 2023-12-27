-- name: CreateUser :one
INSERT INTO users(Email,passwd,id,red,created_at) VALUES($1,$2,$3,$4,$5)
RETURNING *;