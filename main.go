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
	handler "github.com/jaydee029/Verses/handler"
	"github.com/jaydee029/Verses/internal/database"
	"github.com/jaydee029/Verses/pubsub"
	"github.com/joho/godotenv"
)

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

	dbcon, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbcon.Close()

	queries := database.New(dbcon)
	rabbitConnString := os.Getenv("RABBITMQ_CONN")
	if rabbitConnString == "" {
		log.Fatalf("Message broker connection string not set")
	}

	conn, err := pubsub.ConnectToBroker(rabbitConnString)
	if err != nil {
		log.Fatalf("Failed to connect to message broker: %v", err)
	}
	defer conn.Close()

	err = pubsub.InitBroker(conn)
	if err != nil {
		log.Fatalf("Failed to initialize message broker: %v", err)
	}
	h := handler.New(0, jwt_secret, os.Getenv("RED_KEY"), queries, dbcon, conn)

	port := os.Getenv("PORT")

	r := chi.NewRouter()
	s := chi.NewRouter()
	t := chi.NewRouter()

	fileconfig := h.Reqcounts(http.StripPrefix("/app", http.FileServer(http.Dir("./index.html"))))
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
	s.Post("/signup", h.CreateUser)
	s.Post("/login", h.UserLogin)
	s.Post("/prose", h.PostProse)
	s.Get("/{username}/prose", h.GetProse)
	s.Get("/prose/{proseId}", h.ProsebyId)
	s.Post("/prose/{proseId}/togglelike", h.ToggleLike)
	s.Get("/timeline", h.Timeline)
	s.Post("/{proseid}/comments", h.PostComment)
	s.Get("/{proseid}/comments", h.Getcomments)
	s.Post("/comments/{commentid}/togglelike", h.ToggCommentLike)
	s.Post("/refresh", h.VerifyRefresh)
	s.Post("/revoke", h.RevokeToken)
	s.Put("/users", h.UpdateUser)
	s.Get("/users/{username}", h.GetUser)
	s.Get("/users", h.GetUsers)
	s.Post("/users/{username}/toggle_follow", h.ToggleFollow)
	s.Delete("/prose/{proseId}", h.DeleteProse)
	s.Get("/notifications", h.Notifications)
	s.Post("/notifications/{notificationid}/mark_as_read", h.ReadNotification)
	s.Post("/notifications/mark_as_read", h.ReadNotifications)
	s.Post("/gold/webhooks", h.Is_red)
	t.Get("/metrics", h.Metrics)

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
