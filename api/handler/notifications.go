package handler

import (
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaydee029/Verses/api/middleware"
	"github.com/jaydee029/Verses/internal/database"
	"github.com/jaydee029/Verses/pubsub"
	"go.uber.org/zap"
)

type Notification struct {
	ID           pgtype.UUID      `json:"id"`
	Userid       pgtype.UUID      `json:"userid"`
	Proseid      pgtype.UUID      `json:"proseid"`
	Actors       []string         `json:"actors"`
	Generated_at pgtype.Timestamp `json:"generated_at"`
	Read         bool             `json:"read"`
	Type         string           `json:"type"`
}

/*
	type notificationclient struct {
		notifications chan Notification
		Userid        pgtype.UUID
	}
*/
func (cfg *Handler) Notifications(w http.ResponseWriter, r *http.Request) {
	// token, err := auth.BearerHeader(r.Header)

	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "error decoding auth header:"+err.Error())
	// 	return
	// }
	// useridstr, err := auth.ValidateToken(token, cfg.Jwtsecret)

	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "error parsing the userid:"+err.Error())
	// 	return
	// }

	useridstr := r.Context().Value(middleware.UserIDKey).(string)

	var userid pgtype.UUID
	err := userid.Scan(useridstr)
	if err != nil {
		cfg.logger.Info("Error setting UUID:", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if content_type, _, err := mime.ParseMediaType(r.Header.Get("Accept")); err == nil && content_type == "text/event-stream" {
		cfg.subscribeTonotifications(w, r.Context(), userid)
		return
	}

	var before pgtype.Timestamp

	beforestr := r.URL.Query().Get("before")
	if beforestr != "" {
		parsedTime, err := time.Parse(time.RFC3339, beforestr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid timestamp format")
			return
		}
		err = before.Scan(parsedTime)
		if err != nil {
			cfg.logger.Info("Error converting timestamp to pgtype:", zap.Error(err))
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		before = pgtype.Timestamp{
			Valid: false,
		}
	}

	limitstr := r.URL.Query().Get("limit")
	if limitstr == "" {
		limitstr = "10"
	}

	limit, err := strconv.ParseInt(limitstr, 10, 32)
	if err != nil {
		cfg.logger.Info("Error converting limit value to int type:", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// err = before.Scan(beforestr)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	notifications, err := cfg.DB.GetNotifications(r.Context(), database.GetNotificationsParams{
		UserID:  userid,
		Column2: before,
		Limit:   int32(limit),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var Notifications []Notification

	for _, k := range notifications {
		Notifications = append(Notifications, Notification{
			ID:           k.ID,
			Userid:       k.UserID,
			Actors:       k.Actors,
			Generated_at: k.GeneratedAt,
			Type:         k.Type,
			Read:         k.Read,
		})
	}

	respondWithJson(w, http.StatusOK, Notifications)
}

func (cfg *Handler) ReadNotification(w http.ResponseWriter, r *http.Request) {

	// token, err := auth.BearerHeader(r.Header)

	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "error decoding auth header:"+err.Error())
	// 	return
	// }
	// useridstr, err := auth.ValidateToken(token, cfg.Jwtsecret)

	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "error parsing the userid:"+err.Error())
	// 	return
	// }

	useridstr := r.Context().Value(middleware.UserIDKey).(string)
	notificationidstr := chi.URLParam(r, "notificationid")

	if notificationidstr == "" {
		respondWithError(w, http.StatusBadRequest, "notification id not provided")
	}

	var userid pgtype.UUID
	err := userid.Scan(useridstr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var notificationid pgtype.UUID
	err = notificationid.Scan(useridstr)
	if err != nil {
		cfg.logger.Info("Error setting UUID:", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = cfg.DB.ReadNotificationSingle(r.Context(), database.ReadNotificationSingleParams{
		UserID: userid,
		ID:     notificationid,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJson(w, http.StatusNoContent, "Notification Read")
}

func (cfg *Handler) ReadNotifications(w http.ResponseWriter, r *http.Request) {

	// token, err := auth.BearerHeader(r.Header)

	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "error decoding auth header:"+err.Error())
	// 	return
	// }
	// useridstr, err := auth.ValidateToken(token, cfg.Jwtsecret)

	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "error parsing the userid:"+err.Error())
	// 	return
	// }

	useridstr := r.Context().Value(middleware.UserIDKey).(string)
	var userid pgtype.UUID
	err := userid.Scan(useridstr)
	if err != nil {
		cfg.logger.Info("Error setting UUID:", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = cfg.DB.ReadNotificationAll(r.Context(), userid)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJson(w, http.StatusNoContent, "Notifications Read")
}

func (cfg *Handler) Broadcastnotifications(n Notification) {

	err := pubsub.Publish(cfg.pubsub, "notifications_direct", "notification_item."+uuid.UUID(n.Userid.Bytes).String(), n)
	if err != nil {
		cfg.logger.Info("error publishing notification item:", zap.Error(err))
		return
	}

}
