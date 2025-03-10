package handler

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaydee029/Verses/internal/database"
	"go.uber.org/zap"
)

func (cfg *handler) FollowNotification(followeeid, followerid pgtype.UUID) {

	tx, err := cfg.DBpool.Begin(context.Background())
	if err != nil {
		cfg.logger.Info("error starting the transaction", zap.Error(err))
		return
	}

	qtx := cfg.DB.WithTx(tx)

	defer func() {
		if tx != nil {
			tx.Rollback(context.Background())
		}
	}()

	user, err := qtx.GetUserbyId(context.Background(), followerid)
	if err != nil {
		cfg.logger.Info("error fetching the follower username", zap.Error(err))
		return
	}

	actors := []string{user.Username}
	var notification Notification

	notified, err := qtx.NotificationActorExists(context.Background(), database.NotificationActorExistsParams{
		UserID:  followeeid,
		Column2: user.Username,
	})
	if err != nil {
		cfg.logger.Info("error while fetching notification for the actor:", zap.Error(err))
		return
	}

	if !notified {
		cfg.logger.Info("notification not found, now creating..")

		generated_at := time.Now().UTC()
		var pgtype_generated_at pgtype.Timestamp
		if err = pgtype_generated_at.Scan(generated_at); err != nil {
			cfg.logger.Info("error while converting timestamp to pgtype:", zap.Error(err))
			return
		}

		nid := uuid.New().String()
		var pgUUID pgtype.UUID

		err = pgUUID.Scan(nid)
		if err != nil {
			cfg.logger.Info("error while converting notification id to pgtype", zap.Error(err))
			return
		}

		err = qtx.InsertNotification(context.Background(), database.InsertNotificationParams{
			ID:          pgUUID,
			UserID:      followeeid,
			Actors:      actors,
			GeneratedAt: pgtype_generated_at,
			Type:        "follow",
		})
		if err != nil {
			log.Println("error while inserting notification:" + err.Error())
			return
		}
		notification.Actors = actors
		notification.ID = pgUUID
		notification.Generated_at = pgtype_generated_at

	} else {
		generated_at := time.Now().UTC()
		var pgtype_generated_at pgtype.Timestamp
		if err = pgtype_generated_at.Scan(generated_at); err != nil {
			cfg.logger.Info("error while converting timestamp to pgtype:", zap.Error(err))
			return
		}

		notificationid, err := qtx.NotificationExists(context.Background(), followeeid)

		if err != nil {
			cfg.logger.Info("error while fetching notification id for the user:", zap.Error(err))
			return
		}

		actors, err = qtx.UpdateNotification(context.Background(), database.UpdateNotificationParams{
			Column1:     user.Username,
			ID:          notificationid,
			GeneratedAt: pgtype_generated_at,
		})

		if err != nil {
			cfg.logger.Info("error while updating notification:", zap.Error(err))
			return
		}
		notification.Actors = actors
		notification.ID = notificationid
		notification.Generated_at = pgtype_generated_at

	}

	if err = tx.Commit(context.Background()); err != nil {
		cfg.logger.Info("error commmiting the transaction:", zap.Error(err))
		return
	}

	tx = nil

	notification.Userid = followeeid
	notification.Type = "follow"

	go cfg.Broadcastnotifications(notification)
}
