-- name: InserinTimeline :exec
INSERT INTO timeline(prose_id,user_id) VALUES($1,$2);

-- name: FetchTimelineItems :many
INSERT INTO timeline(prose_id,user_id) 
SELECT $1, follower_id FROM follows WHERE followee_id=$2
RETURNING id,user_id;