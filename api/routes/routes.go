package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jaydee029/Verses/api/handler"
	"github.com/jaydee029/Verses/api/middleware"
)

func SetupRoutes(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	apiRouter := chi.NewRouter()

	// Public Routes
	apiRouter.Post("/signup", h.CreateUser)
	apiRouter.Post("/login", h.UserLogin)
	apiRouter.Get("/admin/healthz", middleware.Apireadiness)
	apiRouter.Get("/admin/metrics", h.Metrics)
	//ARRAY_APPEND(notifications.actors,$1),
	// Protected Routes
	apiRouter.Group(func(r chi.Router) {
		r.Use(middleware.Authmiddleware(h.Jwtsecret))
		r.Post("/prose", h.PostProse)
		r.Get("/{username}/prose", h.GetProse)
		r.Get("/prose/{proseId}", h.ProsebyId)
		r.Post("/prose/{proseId}/togglelike", h.ToggleLike)
		r.Delete("/prose/{proseId}", h.DeleteProse)
		r.Get("/timeline", h.Timeline)
		r.Post("/{proseid}/comments", h.PostComment)
		r.Get("/{proseid}/comments", h.Getcomments)
		r.Post("/comments/{commentid}/togglelike", h.ToggCommentLike)
		r.Put("/users", h.UpdateUser)
		r.Get("/users/{username}", h.GetUser)
		r.Get("/users", h.GetUsers)
		r.Get("/users/{username}", h.GetUser)
		r.Post("/users/{username}/toggle_follow", h.ToggleFollow)
		r.Get("/notifications", h.Notifications)
		r.Post("/notifications/{notificationid}/mark_as_read", h.ReadNotification)
		r.Post("/notifications/mark_as_read", h.ReadNotifications)
		r.Post("/refresh", h.VerifyRefresh)
		r.Post("/revoke", h.RevokeToken)
	})

	// Mount /api
	r.Mount("/api", apiRouter)

	return r
}
