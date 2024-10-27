-- name: InserinTimeline :exec
INSERT INTO timeline(prose_id,user_id) VALUES($1,$2);

-- name: FetchTimelineItems :many
INSERT INTO timeline(prose_id,user_id) 
SELECT $1, follower_id FROM follows WHERE followee_id=$2
RETURNING id,user_id;

-- name: GetTimeline :many
SELECT tl.id, tl.prose_id, p.body, p.created_at, p.updated_at, p.likes, p.comments , u.username,
CASE WHEN author_id=$1 THEN true ELSE false END AS Mine,
CASE WHEN Likes.user_id IS NOT NULL THEN true ELSE false END AS Liked
FROM timeline AS tl INNER JOIN prose AS p 
ON tl.prose_id=p.id
INNER JOIN users AS u 
ON p.author_id=u.id
LEFT JOIN post_likes AS Likes
ON Likes.user_id=$1 AND Likes.prose_id=p.id
AND 
$2::TIMESTAMP IS NULL OR p.created_at < $2
WHERE tl.user_id=$1
ORDER BY p.created_at DESC, p.id DESC
LIMIT $3;