package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaydee029/Verses/internal/database"
)

func (cfg *apiconfig) FollowNotification(followeeid, followerid pgtype.UUID) {

	tx, err := cfg.DBpool.Begin(context.Background())
	if err != nil {
		log.Println("error starting the transaction")
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
		log.Println("error fetching the follower username")
		return
	}

	IfNotified, err := qtx.NotificationActorExists(context.Background(), database.NotificationActorExistsParams{
		UserID:  followeeid,
		Column2: user.Username,
	})

	if err != nil {
		log.Panicln("error while fetching notification for the actor:" + err.Error())
		return
	}

	if IfNotified {
		return
	}

	notificationid, err := qtx.NotificationExists(context.Background(), followeeid)

	if err != nil && err != sql.ErrNoRows {
		log.Println("error while fetching notification id for the user:" + err.Error())
		return
	}

	actors := []string{user.Username}

	generated_at := time.Now().UTC()
	var pgtype_generated_at pgtype.Timestamp
	if err = pgtype_generated_at.Scan(generated_at); err != nil {
		log.Println("error while converting timestamp to pgtype:" + err.Error())
		return
	}
	var notification Notification

	if err == sql.ErrNoRows {

		nid := uuid.New()

		if err = notificationid.Scan(nid); err != nil {
			log.Println("error while converting notification id to pgtype" + err.Error())
			return
		}

		err = qtx.InsertNotification(context.Background(), database.InsertNotificationParams{
			ID:          notificationid,
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

	} else {

		actors, err = qtx.UpdateNotification(context.Background(), database.UpdateNotificationParams{
			Column1:     user.Username,
			ID:          notificationid,
			GeneratedAt: pgtype_generated_at,
		})
		if err != nil {
			log.Println("error while updating notification:" + err.Error())
			return
		}
		notification.Actors = actors

	}

	if err = tx.Commit(context.Background()); err != nil {
		log.Println("error commmiting the transaction:" + err.Error())
		return
	}
	tx = nil

	notification.ID = notificationid
	notification.Generated_at = pgtype_generated_at
	notification.Userid = followeeid
	notification.Type = "follow"

	//realtime broadcast

}
