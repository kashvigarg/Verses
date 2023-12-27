-- name: CreateUser :one
INSERT INTO users(Email,passwd,id,red,created_at) VALUES($1,$2,$3,$4,$5)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET Email=$1, passwd=$2 WHERE id==$3
RETURNING *;