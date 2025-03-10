package handler

import (
	"context"
	"time"

	"github.com/jaydee029/Verses/utils"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaydee029/Verses/internal/database"
)

func (cfg *handler) notifycommentmentions(c Comment) {

	mentionedusers, err := utils.Mentions(c.Body)
	if err != nil {
		cfg.logger.Info("Error fetching mentions in the comment:", zap.Error(err))
		return
	}

	generated_at := time.Now().UTC()

	var generated_at_pgtype pgtype.Timestamp

	if err = generated_at_pgtype.Scan(generated_at); err != nil {
		cfg.logger.Info("Error generating timestamp:", zap.Error(err))
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
		cfg.logger.Info("Error fetching comment mentions:", zap.Error(err))
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

func (cfg *handler) notifypostmentions(p Prose) {

	mentionedusers, err := utils.Mentions(p.Body)
	if err != nil {
		cfg.logger.Info("Error fetching mentions in the post:", zap.Error(err))
		return
	}

	generated_at := time.Now().UTC()

	var generated_at_pgtype pgtype.Timestamp

	if err = generated_at_pgtype.Scan(generated_at); err != nil {
		cfg.logger.Info("Error converting timestamp to pgtype format:", zap.Error(err))
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
		cfg.logger.Info("Error fecthing post mentions from the database:", zap.Error(err))
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
