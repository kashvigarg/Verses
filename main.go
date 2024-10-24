package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
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
	DBpool           *pgxpool.Pool
}

//go:embed static/*
var staticFiles embed.FS

func main() {
	godotenv.Load(".env")

	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		log.Fatal("JWT secret key not set")
	}

	dbURL := os.Getenv("DB_CONN")
	if dbURL == "" {
		log.Fatal("database connection string not set")
	}
	/*
		dbcon, err := sql.Open("postgres", dbURL)
		if err != nil {
			fmt.Print(err.Error())
		}*/

	dbcon, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbcon.Close()

	queries := database.New(dbcon)

	apicfg := apiconfig{
		fileservercounts: 0,
		jwtsecret:        jwt_secret,
		apiKey:           os.Getenv("GOLD_KEY"),
		DB:               queries,
		DBpool:           dbcon,
	}

	port := os.Getenv("PORT")

	r := chi.NewRouter()
	s := chi.NewRouter()
	t := chi.NewRouter()

	fileconfig := apicfg.reqcounts(http.StripPrefix("/app", http.FileServer(http.Dir("./index.html"))))
	r.Handle("/app", fileconfig)
	r.Handle("/app/*", fileconfig)

	r.Get("/app", func(w http.ResponseWriter, r *http.Request) {
		file, err := staticFiles.Open("static/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		//w.Write(file)
		if _, err := io.Copy(w, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	s.Get("/healthz", apireadiness)
	s.Post("/prose", apicfg.postProse)
	s.Get("/{username}/prose", apicfg.getProse)
	s.Get("/prose/{proseId}", apicfg.ProsebyId)
	s.Post("/prose/{proseId}/togglelike", apicfg.toggleLike)
	s.Get("/timeline", apicfg.timeline)
	s.Post("/users", apicfg.createUser)
	s.Post("/login", apicfg.userLogin)
	s.Post("/refresh", apicfg.verifyRefresh)
	s.Post("/revoke", apicfg.revokeToken)
	s.Put("/users", apicfg.updateUser)
	s.Post("/users/{username}/toggle_follow", apicfg.toggleFollow)
	s.Delete("/prose/{proseId}", apicfg.DeleteProse)
	s.Post("/gold/webhooks", apicfg.is_gold)
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
