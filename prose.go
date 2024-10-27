package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	auth "github.com/jaydee029/Verses/internal/auth"
	"github.com/jaydee029/Verses/internal/database"
)

func (cfg *apiconfig) getProse(w http.ResponseWriter, r *http.Request) {
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

	username := chi.URLParam(r, "username")

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

	posts, err := cfg.DB.GetsProseAll(r.Context(), database.GetsProseAllParams{
		AuthorID: pgUUID,
		Username: username,
		Column3:  before,
		Limit:    int32(limit),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Prose couldn't be fetched")
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

func (cfg *apiconfig) ProsebyId(w http.ResponseWriter, r *http.Request) {
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

	proseidstr := chi.URLParam(r, "proseId")
	var prose_pgUUID pgtype.UUID

	err = prose_pgUUID.Scan(authorid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var pgUUID pgtype.UUID

	err = pgUUID.Scan(proseidstr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	post, err := cfg.DB.GetProseSingle(r.Context(), database.GetProseSingleParams{
		AuthorID: pgUUID,
		ID:       prose_pgUUID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Prose couldn't be fetched")
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

func (cfg *apiconfig) DeleteProse(w http.ResponseWriter, r *http.Request) {
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
	var pgUUID pgtype.UUID

	err = pgUUID.Scan(authorid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	proseidstr := chi.URLParam(r, "proseId")
	var prose_pgUUID pgtype.UUID

	err = prose_pgUUID.Scan(proseidstr)
	if err != nil {
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
