// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Comment struct {
	ID         int32
	ProseID    pgtype.UUID
	UserID     pgtype.UUID
	Body       string
	CreatedAt  pgtype.Timestamp
	LikesCount int32
}

type CommentLike struct {
	CommentID int32
	UserID    pgtype.UUID
}

type Follow struct {
	FollowerID pgtype.UUID
	FolloweeID pgtype.UUID
}

type Notification struct {
	ID          pgtype.UUID
	UserID      pgtype.UUID
	ProseID     pgtype.UUID
	Actors      []string
	GeneratedAt pgtype.Timestamp
	Type        string
	Read        bool
}

type PostLike struct {
	ProseID pgtype.UUID
	UserID  pgtype.UUID
}

type Prose struct {
	ID        pgtype.UUID
	Body      string
	AuthorID  pgtype.UUID
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	Likes     int32
	Comments  int32
}

type Revocation struct {
	Token     []byte
	RevokedAt pgtype.Timestamp
}

type Timeline struct {
	ID      int32
	ProseID pgtype.UUID
	UserID  pgtype.UUID
}

type User struct {
	Name      string
	Email     string
	Passwd    []byte
	Username  string
	ID        pgtype.UUID
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	IsRed     bool
	Followers int32
	Followees int32
}
