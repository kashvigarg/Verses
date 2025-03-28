package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/jaydee029/Verses/api/handler"
	"github.com/jaydee029/Verses/api/routes"
	"github.com/jaydee029/Verses/internal/database"
	"github.com/jaydee029/Verses/pubsub"
	"github.com/joho/godotenv"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	godotenv.Load(".env")

	logger, _ := zap.NewProduction()

	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		logger.Fatal("JWT secret key not set")
	}

	dbURL := os.Getenv("DB_CONN")
	if dbURL == "" {
		logger.Fatal("database connection string not set")
	}

	dbcon, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Fatal("Unable to connect to database:", zap.Error(err))
		os.Exit(1)
	}

	queries := database.New(dbcon)
	rabbitConnString := os.Getenv("RABBITMQ_CONN")
	if rabbitConnString == "" {
		logger.Fatal("Message broker connection string not set")
	}

	conn, err := pubsub.ConnectToBroker(rabbitConnString, logger)
	if err != nil {
		logger.Fatal("Failed to connect to message broker:", zap.Error(err))
	}

	defer func() {
		logger.Sync()
		dbcon.Close()
		conn.Close()
	}()

	err = pubsub.InitBroker(conn)
	if err != nil {
		logger.Fatal("Failed to initialize message broker:", zap.Error(err))
	}
	h := handler.New(0, jwt_secret, os.Getenv("RED_KEY"), queries, dbcon, conn, logger)

	port := os.Getenv("PORT")

	router := routes.SetupRoutes(h)

	staticFS, _ := fs.Sub(staticFiles, "static")
	fileServer := http.FileServer(http.FS(staticFS))
	router.Handle("/app", http.StripPrefix("/app", fileServer))
	router.Handle("/app/*", http.StripPrefix("/app", fileServer))

	//r := chi.NewRouter()
	// s := chi.NewRouter()
	// t := chi.NewRouter()

	// fileconfig := h.Reqcounts(http.StripPrefix("/app", http.FileServer(http.Dir("./index.html"))))
	// r.Handle("/app", fileconfig)
	// r.Handle("/app/*", fileconfig)

	// r.Get("/app", func(w http.ResponseWriter, r *http.Request) {
	// 	file, err := staticFiles.Open("static/index.html")
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	defer file.Close()
	// 	//w.Write(file)
	// 	if _, err := io.Copy(w, file); err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}
	// })

	// t.Get("/admin/healthz", middleware.Apireadiness)
	// t.Get("/admin/metrics", h.Metrics)
	// t.Post("/signup", h.CreateUser)
	// t.Post("/login", h.UserLogin)
	// s.Post("/prose", h.PostProse)
	// s.Get("/{username}/prose", h.GetProse)
	// s.Get("/prose/{proseId}", h.ProsebyId)
	// s.Post("/prose/{proseId}/togglelike", h.ToggleLike)
	// s.Get("/timeline", h.Timeline)
	// s.Post("/{proseid}/comments", h.PostComment)
	// s.Get("/{proseid}/comments", h.Getcomments)
	// s.Post("/comments/{commentid}/togglelike", h.ToggCommentLike)
	// s.Put("/users", h.UpdateUser)
	// s.Get("/users/{username}", h.GetUser)
	// s.Get("/users", h.GetUsers)
	// s.Post("/users/{username}/toggle_follow", h.ToggleFollow)
	// s.Delete("/prose/{proseId}", h.DeleteProse)
	// s.Get("/notifications", h.Notifications)
	// s.Post("/notifications/{notificationid}/mark_as_read", h.ReadNotification)
	// s.Post("/notifications/mark_as_read", h.ReadNotifications)
	// s.Post("/gold/webhooks", h.Is_red)
	//s.Post("/refresh", h.VerifyRefresh)
	// s.Post("/revoke", h.RevokeToken)

	// r.Mount("/api", s)
	// s.With(middleware.Authmiddleware(h.Jwtsecret))
	// r.Mount("/api", t)
	sermux := corsmiddleware(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: sermux,
	}

	log.Printf("The server is live on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
