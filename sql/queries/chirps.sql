-- name: CreateChirp :one
INSERT INTO chirps(id,body,author_id,created_at,updated_at) VALUES($1,$2,$3,$4,$5)
RETURNING *;
/*
-- name: GetChirps :many
SELECT * FROM chirps WHERE chirps.author_id== users.id
ORDER BY chirps.id;

-- name: GetChirp :one
SELECT * FROM chirps WHERE author_id==$1 AND id==$2;
*/