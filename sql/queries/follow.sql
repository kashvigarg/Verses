-- name: GetIdfromUsername :one
SELECT id FROM users WHERE username=$1;

-- name: If_follows :one
SELECT EXISTS (SELECT 1 FROM follows WHERE follower_id=$1 AND followee_id=$2);

-- name: Updatefollower :one
UPDATE users SET followers= followers+1 AND followees=followees+1 WHERE id=$1 RETURNING followers;

-- name: Deletefollower :one
UPDATE users SET followers= followers-1 AND followees= followees-1 WHERE id=$1 RETURNING followers;

-- name: Addfollower :exec
INSERT INTO follows(followee_id,follower_id) VALUES($1,$2);

-- name: Removefollower :exec
DELETE FROM follows WHERE followee_id=$1 AND follower_id=$2;