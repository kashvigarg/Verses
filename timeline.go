package main

import (
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	auth "github.com/jaydee029/Verses/internal/auth"
	"github.com/jaydee029/Verses/internal/database"
)

type timeline_item struct {
	Id     int         `json:"id"`
	Userid pgtype.UUID `json:"userid,omitempty"`
	Post   Prose       `json:"prose"`
}

func (cfg *apiconfig) timeline(w http.ResponseWriter, r *http.Request) {
	token, err := auth.BearerHeader(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	authorid, err := auth.ValidateToken(token, cfg.jwtsecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
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
		limitstr = "20"
	}

	limit, err := strconv.ParseInt(limitstr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var pgUUID pgtype.UUID
	err = pgUUID.Scan(authorid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = before.Scan(beforestr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tl_items, err := cfg.DB.GetTimeline(r.Context(), database.GetTimelineParams{
		AuthorID: pgUUID,
		Column2:  before,
		Limit:    int32(limit),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "timeline items couldn't be fetched")
		return
	}

	var timeline []timeline_item

	for _, k := range tl_items {
		timeline = append(timeline, timeline_item{
			Id: int(k.ID),
			Post: Prose{
				ID:          k.ProseID,
				Username:    k.Username,
				Body:        k.Body,
				Created_at:  k.CreatedAt,
				Updated_at:  k.UpdatedAt,
				Mine:        k.Mine,
				Liked:       k.Liked,
				Likes_count: int(k.Likes),
				Comments:    int(k.Comments),
			},
		})

		respondWithJson(w, http.StatusOK, timeline)

	}
}
