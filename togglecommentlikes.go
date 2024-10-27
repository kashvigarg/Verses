package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	auth "github.com/jaydee029/Verses/internal/auth"
	"github.com/jaydee029/Verses/internal/database"
)

type toggCommentLike struct {
	Liked       bool `json:"liked"`
	Likes_count int  `json:"likes_count"`
}

func (cfg *apiconfig) toggCommentLike(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user_id, err := auth.ValidateToken(token, cfg.jwtsecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	commentidstr := chi.URLParam(r, "commentid")
	Commentid, err := strconv.Atoi(commentidstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	var pgUUID pgtype.UUID

	err = pgUUID.Scan(user_id)
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

	If_liked, err := qtx.IfCommentLiked(r.Context(), database.IfCommentLikedParams{
		CommentID: int32(Commentid),
		UserID:    pgUUID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var liked bool
	var likes int32

	if If_liked {
		err = qtx.RemoveCommentLike(r.Context(), database.RemoveCommentLikeParams{
			CommentID: int32(Commentid),
			UserID:    pgUUID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error while disliking the comment")
			return
		}

		likes, err = qtx.DecreaseCommentLikeCount(r.Context(), int32(Commentid))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error while decreasing the like count")
			return
		}

		liked = false

	} else {
		err = qtx.AddCommentLike(r.Context(), database.AddCommentLikeParams{
			CommentID: int32(Commentid),
			UserID:    pgUUID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error while liking the comment")
			return
		}

		likes, err = qtx.IncreaseCommentLikeCount(r.Context(), int32(Commentid))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error while increasing the like count")
			return
		}
		liked = true
	}
	tx.Commit(r.Context())
	tx = nil

	respondWithJson(w, http.StatusAccepted, toggCommentLike{
		Liked:       liked,
		Likes_count: int(likes),
	})

}
