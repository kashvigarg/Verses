-- name: CreateUser :one
INSERT INTO users(name,Email,passwd,id,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetUser :one
SELECT passwd,id FROM users WHERE Email==$1;