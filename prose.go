package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	auth "github.com/jaydee029/Verses/internal/auth"
	"github.com/jaydee029/Verses/internal/database"
)

type Prose struct {
	Id         int              `json:"id"`
	Body       string           `json:"body"`
	Created_at pgtype.Timestamp `json:"created_at"`
	Updated_at pgtype.Timestamp `json:"updated_at"`
}

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

	var pgUUID pgtype.UUID

	err = pgUUID.Scan(authorid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirps, err := cfg.DB.GetsProse(r.Context(), pgUUID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Prose couldn't be fetched")
		return
	}
	/*
		author_id := -1
		s := r.URL.Query().Get("author_id")
		if s != "" {
			author_id, err = strconv.Atoi(s)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "author id couldnt be parsed")
				return
			}
		}

			sorting := "asc"
			sort_val := r.URL.Query().Get("sort")

			if sort_val == "desc" {
				sorting = "desc"
			}

			chirps_ := []Chirp{}
			for _, chirp := range chirps {
				if author_id != -1 && chirp.Author_id != author_id {
					continue
				}

				chirps_ = append(chirps_, Chirp{
					Id:        chirp.Id,
					Body:      chirp.Body,
					Author_id: chirp.Author_id,
				})
			}

			sort.Slice(chirps_, func(i, j int) bool {
				if sorting == "desc" {
					return chirps_[i].Id > chirps_[j].Id
				}
				return chirps_[i].Id < chirps_[j].Id
			})
	*/
	respondWithJson(w, http.StatusOK, chirps)
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

	chirpidstr := chi.URLParam(r, "proseId")
	chirpid, err := strconv.Atoi(chirpidstr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't parse id")
		return
	}
	var pgUUID pgtype.UUID

	err = pgUUID.Scan(authorid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirp, err := cfg.DB.Getprose(r.Context(), database.GetproseParams{
		AuthorID: pgUUID,
		ID:       int32(chirpid),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Prose couldn't be fetched")
		return
	}

	respondWithJson(w, http.StatusOK, Prose{
		Id:         int(chirp.ID),
		Body:       chirp.Body,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
	})
}
