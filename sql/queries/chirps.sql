-- name: CreateChirp :one
INSERT INTO chirps(id,body,author_id,created_at,upadated_at) VALUES($1,$2,$3,$4,$5)
RETURNING *;