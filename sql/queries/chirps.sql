-- name: Createchirp :one
INSERT INTO chirps(id,body,author_id,created_at,updated_at) VALUES($1,$2,$3,$4,$5)
RETURNING *;

-- name: Countchirps :one
SELECT COUNT(*) FROM chirps WHERE author_id==$1;

-- name: GetChirp :one
SELECT * FROM chirps WHERE author_id==$1 AND id==$2;

/*
-- name: GetChirps :many
SELECT * FROM chirps WHERE chirps.author_id== users.id
ORDER BY chirps.id;


*/