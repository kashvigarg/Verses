package main

import "github.com/jackc/pgx/v5/pgtype"

type timeline_item struct {
	Id     int         `json:"id"`
	Userid pgtype.UUID `json:"userid"`
	Post   Post        `json:"post"`
}
