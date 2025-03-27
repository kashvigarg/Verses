package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaydee029/Verses/internal/database"
	"go.uber.org/zap"
)

type togglelike struct {
	Liked       bool `json:"liked"`
	Likes_count int  `json:"likes_count"`
}

func (cfg *Handler) ToggleLike(w http.ResponseWriter, r *http.Request) {

	proseidstr := chi.URLParam(r, "proseId")
	// token, err := auth.BearerHeader(r.Header)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, err.Error())
	// 	return
	// }

	// useridstr, err := auth.ValidateToken(token, cfg.Jwtsecret)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, err.Error())
	// 	return
	// }
	useridstr := r.Context().Value("authorid").(string)

	var user_id pgtype.UUID

	err := user_id.Scan(useridstr)
	if err != nil {
		cfg.logger.Info("Error converting Id to pgtype format:", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var prose_id pgtype.UUID

	err = prose_id.Scan(proseidstr)
	if err != nil {
		cfg.logger.Info("Error converting prose Id to pgtype format:", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tx, err := cfg.DBpool.Begin(r.Context())
	if err != nil {
		cfg.logger.Info("Error starting the transaction:", zap.Error(err))
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

	if err = tx.Commit(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error commmiting the transaction:"+err.Error())
	}
	tx = nil

	respondWithJson(w, http.StatusAccepted, togglelike{
		Likes_count: int(likes_count),
		Liked:       liked,
	})
}
