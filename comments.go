package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	auth "github.com/jaydee029/Verses/internal/auth"
	"github.com/jaydee029/Verses/internal/database"
)

type Comment struct {
	Id          int32            `json:"id"`
	Userid      pgtype.UUID      `json:"-"`
	Username    string           `json:"username"`
	Proseid     pgtype.UUID      `json:"proseid"`
	Created_at  pgtype.Timestamp `json:"created_at"`
	Likes_count int              `json:"likes,omitempty"`
	Liked       bool             `json:"liked,omitempty"`
	Mine        bool             `json:"mine,omitempty"`
	Body        string           `json:"body"`
}

func (cfg *apiconfig) postComment(w http.ResponseWriter, r *http.Request) {

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

	decoder := json.NewDecoder(r.Body)
	params := body{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "parameters couldn't be decoded")
		return
	}

	proseidstr := chi.URLParam(r, "proseid")

	cleanText := strings.TrimSpace(params.Body)

	if len([]rune(cleanText)) > 280 || cleanText == "" {
		respondWithError(w, http.StatusBadRequest, "Comment is invalid")
		return
	}

	var pgUUID pgtype.UUID
	err = pgUUID.Scan(authorid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var proseid pgtype.UUID
	err = proseid.Scan(proseidstr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tx, err := cfg.DBpool.Begin(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	qtx := cfg.DB.WithTx(tx)

	defer func() {
		if tx != nil {
			tx.Rollback(r.Context())
		}
	}()

	comment, err := qtx.CreateComment(r.Context(), database.CreateCommentParams{
		ProseID: proseid,
		UserID:  pgUUID,
		Body:    cleanText,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "comment couldn't be created")
		return
	}

	err = qtx.UpdateCommentCount(r.Context(), proseid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "comment count couldn't be updated")
		return
	}

	err = tx.Commit(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't commit the transaction")
	}
	tx = nil

	respondWithJson(w, http.StatusAccepted, Comment{
		Id:         comment.ID,
		Body:       comment.Body,
		Proseid:    comment.ProseID,
		Created_at: comment.CreatedAt,
	})
}

func (cfg *apiconfig) Getcomments(w http.ResponseWriter, r *http.Request) {

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

	before := 0
	beforestr := r.URL.Query().Get("before")

	if beforestr != "" {
		before, err = strconv.Atoi(beforestr)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	limitstr := r.URL.Query().Get("limit")
	if limitstr == "" {
		limitstr = "10"
	}
	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	proseidstr := chi.URLParam(r, "proseid")

	var pgUUID pgtype.UUID
	err = pgUUID.Scan(authorid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var proseid pgtype.UUID
	err = proseid.Scan(proseidstr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	comments, err := cfg.DB.GetComments(r.Context(), database.GetCommentsParams{
		UserID:  pgUUID,
		ProseID: proseid,
		Column3: int32(before),
		Limit:   int32(limit),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Comments couldn't be fetched")
	}

	var Comments []Comment

	for _, k := range comments {
		Comments = append(Comments, Comment{
			Id:          k.ID,
			Username:    k.Username,
			Liked:       k.Liked,
			Mine:        k.Mine,
			Created_at:  k.CreatedAt,
			Likes_count: int(k.LikesCount),
		})
	}

	respondWithJson(w, http.StatusOK, Comments)

}
