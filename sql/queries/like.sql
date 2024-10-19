-- name: If_likes :one
SELECT EXISTS (SELECT 1 FROM post_likes WHERE prose_id=$1 AND user_id=$2);

-- name: Increaselikescount :one
UPDATE prose SET likes= likes+1 WHERE id=$1 RETURNING likes;

-- name: Deletelikescount :one
UPDATE prose SET likes= likes-1 WHERE id=$1 RETURNING likes;

-- name: Deletelike :exec
DELETE FROM post_likes WHERE prose_id=$1 AND user_id=$2;

-- name: Addlike :exec
INSERT INTO post_likes(prose_id,user_id)VALUES ($1,$2);