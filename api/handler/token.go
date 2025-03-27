package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	auth "github.com/jaydee029/Verses/internal/auth"
	"github.com/jaydee029/Verses/internal/database"
	"go.uber.org/zap"
)

type Token struct {
	Token string `json:"token"`
}

func (cfg *Handler) RevokeToken(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "bytes couldn't be converted")
		return
	}

	var pgtime pgtype.Timestamp

	err = pgtime.Scan(time.Now().UTC())
	if err != nil {
		cfg.logger.Info("Error setting timestamp:", zap.Error(err))
	}

	err = cfg.DB.RevokeToken(r.Context(), database.RevokeTokenParams{
		Token:     []byte(token),
		RevokedAt: pgtime,
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	respondWithJson(w, http.StatusOK, "Token Revoked")
}

func (cfg *Handler) VerifyRefresh(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	is_refresh, err := auth.VerifyRefresh(token, cfg.Jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if !is_refresh {
		respondWithError(w, http.StatusUnauthorized, "Header doesn't contain refresh token")
		return
	}

	is_revoked, err := cfg.DB.VerifyRefresh(r.Context(), []byte(token))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if is_revoked {
		respondWithError(w, http.StatusUnauthorized, "Refresh Token revoked")
		return
	}
	Idstr, err := auth.ValidateToken(token, cfg.Jwtsecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	Id, err := uuid.Parse(Idstr)
	if err != nil {
		cfg.logger.Info("Error parsing string to UUID:", zap.Error(err))
	}

	auth_token, err := auth.Tokenize(Id, cfg.Jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, Token{
		Token: auth_token,
	})
}

func (cfg *Handler) Is_red(w http.ResponseWriter, r *http.Request) {
	type user_struct struct {
		User_id pgtype.UUID `json:"user_id"`
	}
	type body struct {
		Event string      `json:"event"`
		Data  user_struct `json:"data"`
	}

	key, err := auth.VerifyAPIkey(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if key != cfg.apiKey {
		respondWithError(w, http.StatusUnauthorized, "Incorrect API Key")
	}

	decoder := json.NewDecoder(r.Body)
	params := body{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	if params.Event == "user.upgraded" {
		user_res, err := cfg.DB.Is_red(r.Context(), true)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJson(w, http.StatusOK, User{
			Name:   user_res.Name,
			Email:  user_res.Email,
			Is_red: user_res.IsRed,
			ID:     params.Data.User_id,
		})
	}

	respondWithJson(w, http.StatusOK, "http request accepted in the webhook")
}
