-- name: Createchirp :one
INSERT INTO chirps(id, body, author_id) VALUES($1,$2,$3)
RETURNING *;