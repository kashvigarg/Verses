package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jaydee029/Verses/internal/auth"
)

type contextKey string

const UserIDKey contextKey = "userID"

//type UserIDKey struct{}

type errresponse struct {
	Error string `json:"error"`
}

func middlewareErrorresponse(w http.ResponseWriter, code int, res interface{}) {
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

func Authmiddleware(tokensecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token, err := auth.BearerHeader(r.Header)
			if err != nil {
				middlewareErrorresponse(w, http.StatusUnauthorized, errresponse{
					Error: err.Error(),
				})
				return
			}

			authorid, err := auth.ValidateToken(token, tokensecret)
			if err != nil {
				middlewareErrorresponse(w, http.StatusUnauthorized, errresponse{
					Error: err.Error(),
				})
				return
			}

			// if authorid == "" {
			// 	middlewareErrorresponse(w, http.StatusUnauthorized, errresponse{
			// 		Error: "Invalid user ID",
			// 	})
			// 	return
			// }

			ctx := context.WithValue(r.Context(), UserIDKey, authorid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}
