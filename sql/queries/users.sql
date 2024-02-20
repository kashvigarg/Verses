-- name: CreateUser :one
INSERT INTO users(name,Email,passwd,id,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET name=$2 ,Email=$3 ,passwd=$4 ,updated_at=$5 WHERE id==$1
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE Email==$1;

-- name: Is_red :one
INSERT INTO users(is_red) VALUES($1)
RETURNING *;