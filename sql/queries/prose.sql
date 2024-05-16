-- name: Createprose :one
INSERT INTO prose(id,body,author_id,created_at,updated_at) VALUES($1,$2,$3,$4,$5)
RETURNING *;

-- name: Countprose :one
SELECT COUNT(*) FROM prose WHERE author_id=$1;

-- name: Getprose :one
SELECT * FROM prose WHERE author_id=$1 AND id=$2;

-- name: GetsProse :many
SELECT id,body,created_at,updated_at FROM prose WHERE author_id=$1
ORDER BY id;

-- name: Deleteprose :exec
DELETE FROM prose WHERE author_id=$1 AND id=$2;