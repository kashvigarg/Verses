-- name: CreateUser :one
INSERT INTO users(name,Email,passwd,id,created_at,updated_at, username) VALUES($1,$2,$3,$4,$5,$6, $7)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET name=$2 ,passwd=$3 ,updated_at=$4 WHERE id=$1
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE Email=$1;

-- name: GetUserbyId :one
SELECT * FROM users WHERE id=$1;

-- name: Is_red :one
INSERT INTO users(is_red) VALUES($1)
RETURNING *;

-- name: Is_Email :one
SELECT EXISTS (SELECT 1 FROM users WHERE Email=$1);

-- name: Is_Username :one
SELECT EXISTS (SELECT 1 FROM users WHERE username=$1);

