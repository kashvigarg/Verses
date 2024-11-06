-- name: MentionPostNotifications :many
INSERT INTO notifications(user_id,type,actors,prose_id,generated_at)
SELECT users.id, "post_mention",$1, $2, $3 FROM users 
WHERE username= ANY($4::VARCHAR[])
AND users.id=$5
RETURNING id,user_id,generated_at;

-- name: MentionCommentNotifications :many
INSERT INTO notifications(user_id,type,actors,prose_id,generated_at)
SELECT users.id, "comment_mention",$1, $2, $3 FROM users 
WHERE username= ANY($4::VARCHAR[])
AND users.id=$5
ON CONFLICT (user_id,prose_id,type,read) DO UPDATE
SET actors=append_array(actors,$1) 
AND generated_at=$3
RETURNING id,user_id,actors,generated_at;