package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"quikchat/internal/adapter/handler/http/handler"
	"quikchat/pkg/middleware"
)

func New(userHandler *handler.UserHandler, friendHandler *handler.FriendHandler, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestLogger)

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
			r.Post("/login", userHandler.Login)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtSecret))

			r.Get("/me", userHandler.GetCurrentUser)

			// User & Profile routes
			r.Get("/users/{username}", userHandler.GetProfileByUsername)
			r.Put("/users/me/profile", userHandler.UpdateProfile)
			r.Post("/users/{username}/block", userHandler.BlockUser)
			r.Delete("/users/{username}/unblock", userHandler.UnblockUser)

			// Friend routes
			r.Route("/friends", func(r chi.Router) {
				r.Get("/", friendHandler.ListFriends)
				r.Delete("/{username}", friendHandler.Unfriend)
				r.Post("/requests", friendHandler.SendRequest)
				r.Get("/requests", friendHandler.ListIncomingRequests)
				r.Put("/requests/{requestID}", friendHandler.RespondToRequest)
			})
		})
	})

	return r
}

