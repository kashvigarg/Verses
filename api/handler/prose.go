package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	auth "github.com/jaydee029/Verses/internal/auth"
	"github.com/jaydee029/Verses/internal/database"
	"go.uber.org/zap"
)

func (cfg *handler) GetProse(w http.ResponseWriter, r *http.Request) {
	// token, err := auth.BearerHeader(r.Header)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, err.Error())
	// 	return
	// }

	// authorid, err := auth.ValidateToken(token, cfg.Jwtsecret)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, err.Error())
	// 	return
	// }

	authorid := r.Context().Value("authorid").(string)

	username := chi.URLParam(r, "username")

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
	}

	var pgUUID pgtype.UUID
	err = pgUUID.Scan(authorid)
	if err != nil {
		cfg.logger.Info("Error setting UUID:", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//string(pgUUID.Bytes)
	posts, err := cfg.DB.GetsProseAll(r.Context(), database.GetsProseAllParams{
		AuthorID: pgUUID,
		Username: username,
		Column3:  before,
		Limit:    int32(limit),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Prose couldn't be fetched:"+err.Error())
		return
	}

	var prose []Prose

	for _, k := range posts {
		prose = append(prose, Prose{
			ID:          k.ID,
			Body:        k.Body,
			Mine:        k.Mine,
			Liked:       k.Liked,
			Likes_count: int(k.Likes),
			Created_at:  k.CreatedAt,
			Updated_at:  k.UpdatedAt,
			Comments:    int(k.Comments),
		})
	}

	respondWithJson(w, http.StatusOK, prose)
}

func (cfg *handler) ProsebyId(w http.ResponseWriter, r *http.Request) {
	// token, err := auth.BearerHeader(r.Header)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, err.Error())
	// 	return
	// }

	// authorid, err := auth.ValidateToken(token, cfg.Jwtsecret)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, err.Error())
	// 	return
	// }

	authorid := r.Context().Value("authorid").(string)
	proseidstr := chi.URLParam(r, "proseId")
	var prose_pgUUID pgtype.UUID

	err := prose_pgUUID.Scan(proseidstr)
	if err != nil {
		cfg.logger.Info("Error converting Id to pgtype format:", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var pgUUID pgtype.UUID

	err = pgUUID.Scan(authorid)
	if err != nil {
		cfg.logger.Info("Error setting UUID:", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	post, err := cfg.DB.GetProseSingle(r.Context(), database.GetProseSingleParams{
		AuthorID: pgUUID,
		ID:       prose_pgUUID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Prose couldn't be fetched:"+err.Error())
		return
	}

	prose := Prose{
		Username:    post.Username,
		Body:        post.Body,
		Created_at:  post.CreatedAt,
		Updated_at:  post.UpdatedAt,
		Mine:        post.Mine,
		Liked:       post.Liked,
		Likes_count: int(post.Likes),
		Comments:    int(post.Comments),
	}

	respondWithJson(w, http.StatusOK, prose)
}

func (cfg *handler) DeleteProse(w http.ResponseWriter, r *http.Request) {
	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	authorid, err := auth.ValidateToken(token, cfg.Jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	var pgUUID pgtype.UUID

	err = pgUUID.Scan(authorid)
	if err != nil {
		cfg.logger.Info("Error setting UUID:", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	proseidstr := chi.URLParam(r, "proseId")
	var prose_pgUUID pgtype.UUID

	err = prose_pgUUID.Scan(proseidstr)
	if err != nil {
		cfg.logger.Info("Error converting Id to pgtype format:", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = cfg.DB.Deleteprose(r.Context(), database.DeleteproseParams{
		AuthorID: pgUUID,
		ID:       prose_pgUUID,
	})

	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, "Prose deleted")
}
