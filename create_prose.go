package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	auth "github.com/jaydee029/Verses/internal"
	"github.com/jaydee029/Verses/internal/database"
)

type Res struct {
	ID         int              `json:"id"`
	Body       string           `json:"body"`
	Created_at pgtype.Timestamp `json:"created_at"`
	Updated_at pgtype.Timestamp `json:"Updated_at"`
}

type body struct {
	Body string `json:"body"`
}

func (cfg *apiconfig) postProse(w http.ResponseWriter, r *http.Request) {

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

	/*authorid_num, err := uuid.Parse(authorid)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "bytes couldn't be converted")
		return
	}*/

	decoder := json.NewDecoder(r.Body)
	params := body{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "parameters couldn't be decoded")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Prose is too long")
		return
	}
	var pgUUID pgtype.UUID

	err = pgUUID.Scan(authorid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	total, _ := cfg.DB.Countprose(r.Context(), pgUUID)
	content := profane(params.Body)

	var pgtime pgtype.Timestamp

	err = pgtime.Scan(time.Now().UTC())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirp, err := cfg.DB.Createprose(r.Context(), database.CreateproseParams{
		ID:        int32(total + 1),
		AuthorID:  pgUUID,
		Body:      content,
		CreatedAt: pgtime,
		UpdatedAt: pgtime,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create Prose")
		return
	}

	respondWithJson(w, http.StatusCreated, Res{
		ID:         int(chirp.ID),
		Body:       chirp.Body,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
	})
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
	/*authorid_num, err := uuid.Parse(authorid)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "bytes couldn't be converted")
		return
	}*/

	chirpidstr := chi.URLParam(r, "proseId")
	chirpid, err := strconv.Atoi(chirpidstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't parse prose id")
		return
	}

	err = cfg.DB.Deleteprose(r.Context(), database.DeleteproseParams{
		AuthorID: pgUUID,
		ID:       int32(chirpid),
	})

	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, "Prose deleted")
}

func profane(content string) string {
	contentslice := strings.Split(content, " ")

	for i, word := range contentslice {
		wordl := strings.ToLower(word)
		if wordl == "fuck" || wordl == "shit" || wordl == "fornax" {
			contentslice[i] = "****"
		}
	}

	return strings.Join(contentslice, " ")
}
