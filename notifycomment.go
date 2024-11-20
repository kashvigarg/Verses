package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaydee029/Verses/internal/database"
)

func (cfg *apiconfig) CommentNotification(c Comment) {

	nid := uuid.New()
	var nid_pgtype pgtype.UUID
	if err := nid_pgtype.Scan(nid); err != nil {
		log.Println("error while converting notification id to pgtype:" + err.Error())
		return
	}

	generated_at := time.Now().UTC()
	var generated_at_pgtype pgtype.Timestamp
	if err := generated_at_pgtype.Scan(generated_at); err != nil {
		log.Println("error while converting timestamp to pgtype:" + err.Error())
		return
	}

	notification, err := cfg.DB.InsertCommentNotification(context.Background(), database.InsertCommentNotificationParams{
		UserID:      c.Userid,
		ProseID:     c.Proseid,
		GeneratedAt: generated_at_pgtype,
		ID:          nid_pgtype,
		Actors:      []string{c.Username},
	})

	if err != nil {
		log.Println("error while inserting comment notifications:" + err.Error())
		return
	}
	var n Notification

	n.ID = nid_pgtype
	n.Actors = notification.Actors
	n.Userid = c.Userid
	n.Proseid = c.Proseid
	n.Generated_at = notification.GeneratedAt
	n.Type = "comment"

	go cfg.Broadcastnotifications(n)

}
