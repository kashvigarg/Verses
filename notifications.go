package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	auth "github.com/jaydee029/Verses/internal/auth"
	"github.com/jaydee029/Verses/internal/database"
)

type Notification struct {
	ID           pgtype.UUID      `json:"id"`
	User_id      pgtype.UUID      `json:"user_id"`
	Actors       []string         `json:"actors"`
	Generated_at pgtype.Timestamp `json:"generated_at"`
	Read         bool             `json:"read"`
	Type         string           `json:"type"`
}

func (cfg *apiconfig) Notifications(w http.ResponseWriter, r *http.Request) {
	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error decoding auth header:"+err.Error())
		return
	}
	useridstr, err := auth.ValidateToken(token, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error parsing the userid:"+err.Error())
		return
	}

	var before pgtype.Timestamp

	beforestr := r.URL.Query().Get("before")
	if beforestr != "" {
		err = before.Scan(beforestr)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
	}
	limitstr := r.URL.Query().Get("limit")
	if limitstr == "" {
		limitstr = "10"
	}

	limit, err := strconv.ParseInt(limitstr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var userid pgtype.UUID
	err = userid.Scan(useridstr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = before.Scan(beforestr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	notifications, err := cfg.DB.GetNotifications(r.Context(), database.GetNotificationsParams{
		UserID:  userid,
		Column2: before,
		Limit:   int32(limit),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var Notifications []Notification

	for _, k := range notifications {
		Notifications = append(Notifications, Notification{
			ID:           k.ID,
			User_id:      k.UserID,
			Actors:       k.Actors,
			Generated_at: k.GeneratedAt,
			Type:         k.Type,
			Read:         k.Read,
		})
	}

	respondWithJson(w, http.StatusOK, Notifications)
}

func (cfg *apiconfig) ReadNotification(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error decoding auth header:"+err.Error())
		return
	}
	useridstr, err := auth.ValidateToken(token, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error parsing the userid:"+err.Error())
		return
	}

	notificationidstr := chi.URLParam(r, "notificationid")

	if notificationidstr == "" {
		respondWithError(w, http.StatusBadRequest, "notification id not provided")
	}

	var userid pgtype.UUID
	err = userid.Scan(useridstr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var notificationid pgtype.UUID
	err = notificationid.Scan(useridstr)
	if err != nil {
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

func (cfg *apiconfig) ReadNotifications(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error decoding auth header:"+err.Error())
		return
	}
	useridstr, err := auth.ValidateToken(token, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error parsing the userid:"+err.Error())
		return
	}

	var userid pgtype.UUID
	err = userid.Scan(useridstr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = cfg.DB.ReadNotificationAll(r.Context(), userid)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJson(w, http.StatusNoContent, "Notifications Read")
}
