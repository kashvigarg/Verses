package main

import (
	"context"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaydee029/Verses/internal/database"
)

func (cfg *apiconfig) notifycommentmentions(c Comment) {

	mentionedusers, err := mentions(c.Body)
	if err != nil {
		log.Println(err)
		return
	}

	generated_at := time.Now().UTC()

	var generated_at_pgtype pgtype.Timestamp

	if err = generated_at_pgtype.Scan(generated_at); err != nil {
		log.Println(err)
		return
	}

	notifications, err := cfg.DB.MentionCommentNotifications(context.Background(), database.MentionCommentNotificationsParams{
		Actors:      []string{c.User.Username},
		ProseID:     c.Proseid,
		GeneratedAt: generated_at_pgtype,
		Column4:     mentionedusers,
		ID:          c.Userid,
	})

	if err != nil {
		log.Println(err)
		return
	}

	for _, k := range notifications {
		var n Notification

		n.ID = k.ID
		n.Userid = k.UserID
		n.Generated_at = k.GeneratedAt
		n.Actors = k.Actors
		n.Type = "comment_mention"

		go cfg.Broadcastnotifications(n)

	}
}

/*
type Notification struct {
	ID           pgtype.UUID      `json:"id"`
	Userid       pgtype.UUID      `json:"userid"`
	Proseid      pgtype.UUID      `json:"proseid"`
	Actors       []string         `json:"actors"`
	Generated_at pgtype.Timestamp `json:"generated_at"`
	Read         bool             `json:"read"`
	Type         string           `json:"type"`
}
*/

func (cfg *apiconfig) notifypostmentions(p Prose) {

	mentionedusers, err := mentions(p.Body)
	if err != nil {
		log.Println(err)
		return
	}

	generated_at := time.Now().UTC()

	var generated_at_pgtype pgtype.Timestamp

	if err = generated_at_pgtype.Scan(generated_at); err != nil {
		log.Println(err)
		return
	}

	notifications, err := cfg.DB.MentionPostNotifications(context.Background(), database.MentionPostNotificationsParams{
		Actors:      []string{p.User.Username},
		ProseID:     p.ID,
		GeneratedAt: generated_at_pgtype,
		Column4:     mentionedusers,
		ID:          p.Userid,
	})

	if err != nil {
		log.Println(err)
		return
	}

	for _, k := range notifications {
		var n Notification

		n.ID = k.ID
		n.Userid = k.UserID
		n.Generated_at = k.GeneratedAt
		n.Actors = []string{p.User.Username}
		n.Type = "post_mention"

		go cfg.Broadcastnotifications(n)

	}
}

func mentions(content string) ([]string, error) {

	words := strings.Split(content, " ")
	uniqueusers := make(map[string]bool)
	users := []string{}

	pattern := regexp.MustCompile(`^@[a-zA-Z_][a-zA-Z0-9._%+-]{0,8}$`)

	for _, k := range words {

		if pattern.MatchString(k) {
			username := k[1:]

			if !uniqueusers[username] {
				if _, ok := uniqueusers[k[1:]]; !ok {
					uniqueusers[username] = true
					users = append(users, username)
				}
			}
		}
	}

	return users, nil
}
