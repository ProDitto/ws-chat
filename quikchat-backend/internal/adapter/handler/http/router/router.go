package router

import (
	"net/http"
	"quikchat/internal/adapter/handler/http/handler"
	"quikchat/internal/adapter/ws"
	"quikchat/pkg/middleware"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

func New(userHandler *handler.UserHandler, friendHandler *handler.FriendHandler, messageHandler *handler.MessageHandler, hub *ws.Hub, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	r.Use(chi_middleware.Recoverer)
	r.Use(middleware.RequestLogger)

	// Serve static frontend files
	fs := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fs)

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
			r.Post("/login", userHandler.Login)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtSecret))

			// User routes
			r.Get("/users/me", userHandler.GetCurrentUser)
			r.Get("/users/{username}", userHandler.GetProfileByUsername)
			r.Put("/users/me/profile", userHandler.UpdateProfile)
			r.Post("/users/block/{username}", userHandler.BlockUser)
			r.Delete("/users/unblock/{username}", userHandler.UnblockUser)

			// Friend routes
			r.Post("/friends/requests", friendHandler.SendRequest)
			r.Get("/friends/requests/incoming", friendHandler.ListIncomingRequests)
			r.Post("/friends/requests/{requestID}", friendHandler.RespondToRequest)
			r.Get("/friends", friendHandler.ListFriends)
			r.Delete("/friends/{username}", friendHandler.Unfriend)

			// Message Routes
			r.Get("/conversations/{conversationID}/messages", messageHandler.GetMessages)
			r.Post("/conversations/{conversationID}/messages", messageHandler.SendMessage)

			// Media Upload URL Route
			r.Get("/media/presigned-url", messageHandler.GetPresignedURL)
		})

		// WebSocket Route (auth middleware checks header or query param)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtSecret))
			r.Get("/ws", messageHandler.ServeWs)
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	return r
}

