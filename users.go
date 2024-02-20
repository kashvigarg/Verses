package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	auth "github.com/jaydee029/Barkin/internal"
	"github.com/jaydee029/Barkin/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type Input struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
type res struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}
type res_login struct {
	Email         string `json:"email"`
	Token         string `json:"token"`
	Refresh_token string `json:"refresh_token"`
}
type User struct {
	Name     string `json:"name"`
	Password []byte `json:"password"`
	Email    string `json:"email"`
}
type Token struct {
	Token string `json:"token"`
}

func (cfg *apiconfig) createUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := Input{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:     params.Email,
		Passwd:    encrypted,
		ID:        uuid.UUID{},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	respondWithJson(w, http.StatusCreated, res{
		Email: user.Email,
		ID:    user.ID,
	})
}

func (cfg *apiconfig) userLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := Input{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUser(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "user not found")
		return
	}
	err = bcrypt.CompareHashAndPassword(user.Passwd, []byte(params.Password))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password doesn't match")
	}

	Token, err := auth.Tokenize(user.ID, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	Refresh_token, err := auth.RefreshToken(user.ID, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, res_login{
		Email:         params.Email,
		Token:         Token,
		Refresh_token: Refresh_token,
	})

}

func (cfg *apiconfig) updateUser(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	is_refresh, err := auth.VerifyRefresh(token, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if is_refresh {
		respondWithError(w, http.StatusUnauthorized, "Header contains refresh token")
		return
	}

	Idstr, err := auth.ValidateToken(token, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := uuid.FromBytes([]byte(Idstr))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user Id couldn't be parsed")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := User{}
	err = decoder.Decode(&params)

	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	updateduser, err := cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:        userId,
		Name:      params.Name,
		Email:     params.Email,
		Passwd:    hashedPasswd,
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, res{
		Name:  updateduser.Name,
		Email: updateduser.Email,
	})
}

func (cfg *apiconfig) revokeToken(w http.ResponseWriter, r *http.Request) {
	/*
		decoder := json.NewDecoder(r.Body)
		params := User{}
		err := decoder.Decode(&params)

		if err != io.EOF {
			respondWithError(w, http.StatusUnauthorized, "Body is provided")
			return
		}
	*/
	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "bytes couldn't be converted")
		return
	}

	err = cfg.DB.RevokeToken(r.Context(), database.RevokeTokenParams{
		Token:     []byte(token),
		RevokedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	respondWithJson(w, http.StatusOK, "Token Revoked")
}

func (cfg *apiconfig) verifyRefresh(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	is_refresh, err := auth.VerifyRefresh(token, cfg.jwtsecret)

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
	Idstr, err := auth.ValidateToken(token, cfg.jwtsecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userid, err := uuid.FromBytes([]byte(Idstr))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user Id couldn't be parsed")
		return
	}
	auth_token, err := auth.Tokenize(userid, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, Token{
		Token: auth_token,
	})
}

/*
func (cfg *apiconfig) is_red(w http.ResponseWriter, r *http.Request) {
	type user_struct struct {
		User_id int `json:"user_id"`
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
		user_res, err := cfg.DB.Is_red(params.Data.User_id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJson(w, http.StatusOK, res{
			Email:         user_res.Email,
			Is_chirpy_red: user_res.Is_chirpy_red,
			ID:            params.Data.User_id,
		})
	}

	respondWithJson(w, http.StatusOK, "http request accepted in the webhook")
}
*/
