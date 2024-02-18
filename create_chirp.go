package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	auth "github.com/jaydee029/Barkin/internal"
	"github.com/jaydee029/Barkin/internal/database"
)

type Res struct {
	ID         int       `json:"id"`
	Body       string    `json:"body"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"Updated_at"`
}

type body struct {
	Body string `json:"body"`
}

func (cfg *apiconfig) postChirps(w http.ResponseWriter, r *http.Request) {

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

	authorid_num, err := uuid.FromBytes([]byte(authorid))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "bytes couldn't be converted")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := body{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "parameters couldn't be decoded")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	total, _ := cfg.DB.Countchirps(r.Context(), authorid_num)
	content := profane(params.Body)

	chirp, err := cfg.DB.Createchirp(r.Context(), database.CreatechirpParams{
		ID:        int32(total + 1),
		AuthorID:  authorid_num,
		Body:      content,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create chirp")
		return
	}

	respondWithJson(w, http.StatusCreated, Res{
		ID:         int(chirp.ID),
		Body:       chirp.Body,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
	})
}

func (cfg *apiconfig) DeleteChirps(w http.ResponseWriter, r *http.Request) {
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

	authorid_num, err := uuid.FromBytes([]byte(authorid))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "bytes couldn't be converted")
		return
	}

	chirpidstr := chi.URLParam(r, "chirpId")
	chirpid, err := strconv.Atoi(chirpidstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't parse chirp id")
		return
	}

	err = cfg.DB.DeleteChirp(r.Context(), database.DeleteChirpParams{
		AuthorID: authorid_num,
		ID:       int32(chirpid),
	})

	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, "Chirp deleted")
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

func respondWithError(w http.ResponseWriter, code int, res string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", res)
	}
	type errresponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errresponse{
		Error: res,
	})
}

func respondWithJson(w http.ResponseWriter, code int, res interface{}) {
	w.Header().Set("content-type", "application/json")
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}
