package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"quikchat/internal/adapter/handler/http/handler"
	"quikchat/pkg/middleware"
)

func New(userHandler *handler.UserHandler, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestLogger)

	// Public routes
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
	})

	// Protected routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Auth(jwtSecret))
		r.Get("/me", userHandler.GetCurrentUser)
	})

	return r
}

