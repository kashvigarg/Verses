package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	auth "github.com/jaydee029/Verses/internal/auth"
	"github.com/jaydee029/Verses/internal/database"
)

func (cfg *apiconfig) getUser(w http.ResponseWriter, r *http.Request) {

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

	var pgUUID pgtype.UUID
	err = pgUUID.Scan(authorid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.DB.GetUsersingle(r.Context(), database.GetUsersingleParams{
		FolloweeID: pgUUID,
		Username:   username,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User profile couldn't be fetched")
		return
	}

	respondWithJson(w, http.StatusOK, User{
		Name:      user.Name,
		Username:  user.Username,
		ID:        user.ID,
		Follower:  user.Follower,
		Following: user.Following,
	})
}

func (cfg *apiconfig) getUsers(w http.ResponseWriter, r *http.Request) {

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

	after := r.URL.Query().Get("username")

	var pgUUID pgtype.UUID
	err = pgUUID.Scan(authorid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	limitstr := r.URL.Query().Get("limit")
	if limitstr == "" {
		limitstr = "10"
	}

	limit, err := strconv.ParseInt(limitstr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	users, err := cfg.DB.GetUsers(r.Context(), database.GetUsersParams{
		FolloweeID: pgUUID,
		Username:   after,
		Limit:      int32(limit),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Users couldn't be retrieved")
		return
	}

	var Users []User

	for _, user := range users {
		Users = append(Users, User{
			Name:      user.Name,
			Username:  user.Username,
			Follower:  user.Follower,
			Following: user.Following,
			ID:        user.ID,
		})
	}

	respondWithJson(w, http.StatusOK, Users)
}
