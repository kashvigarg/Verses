package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaydee029/Verses/internal/database"
)

func (cfg *apiconfig) BroadcastComment(c Comment) {

	user, err := cfg.DB.GetUserbyId(context.Background(), c.Userid)
	if err != nil {
		log.Println("error fetching the user from the id:" + err.Error())
		return
	}

	c.User = &User{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		Is_red:   user.IsRed,
	}
	c.Mine = false

	go cfg.CommentNotification(c)
	go cfg.notifycommentmentions(c)
	//broadcast comment
}

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

	_, err := cfg.DB.InsertCommentNotification(context.Background(), database.InsertCommentNotificationParams{
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

}
