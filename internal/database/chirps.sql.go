// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: chirps.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createchirp = `-- name: Createchirp :one
INSERT INTO chirps(id, body, author_id) VALUES($1,$2,$3)
RETURNING id, body, author_id
`

type CreatechirpParams struct {
	ID       sql.NullInt16
	Body     sql.NullString
	AuthorID uuid.NullUUID
}

func (q *Queries) Createchirp(ctx context.Context, arg CreatechirpParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, createchirp, arg.ID, arg.Body, arg.AuthorID)
	var i Chirp
	err := row.Scan(&i.ID, &i.Body, &i.AuthorID)
	return i, err
}
