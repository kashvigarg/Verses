package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	auth "github.com/jaydee029/Verses/internal/auth"
	"github.com/jaydee029/Verses/internal/database"
)

type togglelike struct {
	Liked       bool `json:"liked"`
	Likes_count int  `json:"likes_count"`
}

func (cfg *apiconfig) toggleLike(w http.ResponseWriter, r *http.Request) {

	proseidstr := chi.URLParam(r, "postid")

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	useridstr, err := auth.ValidateToken(token, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var user_id pgtype.UUID

	err = user_id.Scan(useridstr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var prose_id pgtype.UUID

	err = prose_id.Scan(proseidstr)
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

	if_liked, err := qtx.If_likes(r.Context(), database.If_likesParams{
		ProseID: prose_id,
		UserID:  user_id,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching likes")
		return
	}

	var likes_count int32
	var liked bool

	if if_liked {
		err = qtx.Deletelike(r.Context(), database.DeletelikeParams{
			ProseID: prose_id,
			UserID:  user_id,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error deleting like")
			return
		}

		likes_count, err = qtx.Deletelikescount(r.Context(), prose_id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error decreasing likes count")
			return
		}
		liked = false
	} else {
		err = qtx.Addlike(r.Context(), database.AddlikeParams{
			ProseID: prose_id,
			UserID:  user_id,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error adding like")
			return
		}

		likes_count, err = qtx.Increaselikescount(r.Context(), prose_id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error increasing likes count")
			return
		}
		liked = true
	}
	tx.Commit(r.Context())
	tx = nil
	respondWithJson(w, http.StatusAccepted, togglelike{
		Likes_count: int(likes_count),
		Liked:       liked,
	})
}
