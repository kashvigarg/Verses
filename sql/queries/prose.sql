-- name: Createprose :one
INSERT INTO prose(id,body,author_id,created_at,updated_at) VALUES($1,$2,$3,$4,$5)
RETURNING *;

-- name: Countprose :one
SELECT COUNT(*) FROM prose WHERE author_id=$1;

-- name: GetProseSingle :one
SELECT p.body,p.id,p.created_at,p.updated_at,p.likes, p.comments ,u.username ,
CASE WHEN author_id=$1 THEN true ELSE false END AS Mine,
CASE WHEN Likes.user_id IS NOT NULL THEN true ELSE false END AS Liked
FROM prose as p 
INNER JOIN users AS u 
ON p.author_id=u.id
LEFT JOIN post_likes AS Likes
ON Likes.user_id=$1 AND Likes.prose_id=p.id
WHERE p.id=$2;

-- name: GetsProseAll :many
SELECT id,body,created_at,updated_at,likes, comments,
CASE WHEN author_id=$1 THEN true ELSE false END AS Mine, 
CASE WHEN Likes.user_id IS NOT NULL THEN true ELSE false END AS Liked
FROM prose LEFT JOIN post_likes AS Likes 
ON Likes.user_id=$1 AND Likes.prose_id=prose.id 
WHERE prose.author_id=(SELECT id FROM users WHERE username=$2)
AND
$3::TIMESTAMP IS NULL OR prose.created_at < $3
ORDER BY prose.created_at DESC,prose.id DESC
LIMIT $4;

-- name: Deleteprose :exec
DELETE FROM prose WHERE author_id=$1 AND id=$2;