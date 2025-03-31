// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: mentions.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const mentionCommentNotifications = `-- name: MentionCommentNotifications :many
INSERT INTO notifications(user_id,type,actors,prose_id,generated_at)
SELECT users.id, 'comment_mention',$1, $2, $3 FROM users 
WHERE username= ANY($4::VARCHAR[])
AND users.id=$5
ON CONFLICT (user_id,prose_id,type,read) DO UPDATE
SET actors= notifications.actors || $1,
generated_at=$3
RETURNING id,user_id,actors,generated_at
`

type MentionCommentNotificationsParams struct {
	Actors      []string
	ProseID     pgtype.UUID
	GeneratedAt pgtype.Timestamp
	Column4     []string
	ID          pgtype.UUID
}

type MentionCommentNotificationsRow struct {
	ID          pgtype.UUID
	UserID      pgtype.UUID
	Actors      []string
	GeneratedAt pgtype.Timestamp
}

func (q *Queries) MentionCommentNotifications(ctx context.Context, arg MentionCommentNotificationsParams) ([]MentionCommentNotificationsRow, error) {
	rows, err := q.db.Query(ctx, mentionCommentNotifications,
		arg.Actors,
		arg.ProseID,
		arg.GeneratedAt,
		arg.Column4,
		arg.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MentionCommentNotificationsRow
	for rows.Next() {
		var i MentionCommentNotificationsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Actors,
			&i.GeneratedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const mentionPostNotifications = `-- name: MentionPostNotifications :many
INSERT INTO notifications(user_id,type,actors,prose_id,generated_at)
SELECT users.id, 'post_mention',$1, $2, $3 FROM users 
WHERE username= ANY($4::VARCHAR[])
AND users.id=$5
RETURNING id,user_id,generated_at
`

type MentionPostNotificationsParams struct {
	Actors      []string
	ProseID     pgtype.UUID
	GeneratedAt pgtype.Timestamp
	Column4     []string
	ID          pgtype.UUID
}

type MentionPostNotificationsRow struct {
	ID          pgtype.UUID
	UserID      pgtype.UUID
	GeneratedAt pgtype.Timestamp
}

func (q *Queries) MentionPostNotifications(ctx context.Context, arg MentionPostNotificationsParams) ([]MentionPostNotificationsRow, error) {
	rows, err := q.db.Query(ctx, mentionPostNotifications,
		arg.Actors,
		arg.ProseID,
		arg.GeneratedAt,
		arg.Column4,
		arg.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MentionPostNotificationsRow
	for rows.Next() {
		var i MentionPostNotificationsRow
		if err := rows.Scan(&i.ID, &i.UserID, &i.GeneratedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
