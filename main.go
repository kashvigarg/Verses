package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/jaydee029/Verses/internal/database"
	"github.com/joho/godotenv"
)

type apiconfig struct {
	fileservercounts int
	jwtsecret        string
	apiKey           string
	DB               *database.Queries
}

func main() {
	err := godotenv.Load(".env")
	/*
		if err != nil {
			log.Fatal(".env file not provided")
		}
	*/
	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		log.Fatal("JWT secret key not set")
	}

	dbURL := os.Getenv("DB_CONN")
	if dbURL == "" {
		log.Fatal("database connection string not set")
	}

	dbcon, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Print(err.Error())
	}
	queries := database.New(dbcon)

	apicfg := apiconfig{
		fileservercounts: 0,
		jwtsecret:        jwt_secret,
		apiKey:           os.Getenv("POLKA_KEY"),
		DB:               queries,
	}

	port := os.Getenv("PORT")

	r := chi.NewRouter()
	s := chi.NewRouter()
	t := chi.NewRouter()

	fileconfig := apicfg.reqcounts(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app", fileconfig)
	r.Handle("/app/*", fileconfig)

	s.Get("/healthz", apireadiness)
	s.Post("/chirps", apicfg.postProse)
	s.Get("/chirps", apicfg.getProse)
	s.Get("/chirps/{chirpId}", apicfg.ProsebyId)
	s.Post("/users", apicfg.createUser)
	s.Post("/login", apicfg.userLogin)
	s.Post("/refresh", apicfg.verifyRefresh)
	s.Post("/revoke", apicfg.revokeToken)
	s.Put("/users", apicfg.updateUser)
	s.Delete("/chirps/{chirpId}", apicfg.DeleteProse)
	s.Post("/polka/webhooks", apicfg.is_red)
	t.Get("/metrics", apicfg.metrics)

	r.Mount("/api", s)
	r.Mount("/admin", t)
	sermux := corsmiddleware(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: sermux,
	}

	log.Printf("The server is live on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
