-- name: GetIdfromUsername :one
SELECT id FROM users WHERE username=$1;

-- name: If_follows :one
SELECT EXISTS (SELECT 1 FROM follows WHERE follower_id=$1 AND followee_id=$2);

-- name: Updatefollower :one
UPDATE users SET followers= followers+1 WHERE id=$1 RETURNING followers;

-- name: Updatefollowee :exec
UPDATE users SET followees= followees+1 WHERE id=$1;

-- name: Deletefollower :one
UPDATE users SET followers= followers-1 WHERE id=$1 RETURNING followers;

-- name: Deletefollowee :exec
UPDATE users SET followees= followees-1 WHERE id=$1;


