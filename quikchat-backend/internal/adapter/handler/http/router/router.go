package router

import (
	"net/http"
	"quikchat/internal/adapter/handler/http/handler"
	"quikchat/internal/adapter/ws"
	"quikchat/pkg/middleware"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

func New(
	userHandler *handler.UserHandler,
	friendHandler *handler.FriendHandler,
	messageHandler *handler.MessageHandler,
	groupHandler *handler.GroupHandler,
	notificationHandler *handler.NotificationHandler,
	hub *ws.Hub,
	jwtSecret string,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chi_middleware.Recoverer)
	r.Use(middleware.RequestLogger)

	// Public routes
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
	})

	// Protected routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Auth(jwtSecret))

		// User routes
		r.Get("/users/me", userHandler.GetCurrentUser)
		r.Get("/users/{username}", userHandler.GetProfileByUsername)
		r.Put("/users/me/profile", userHandler.UpdateProfile)
		r.Post("/users/{username}/block", userHandler.BlockUser)
		r.Delete("/users/{username}/unblock", userHandler.UnblockUser)

		// Friend routes
		r.Get("/friends", friendHandler.ListFriends)
		r.Post("/friends/requests", friendHandler.SendRequest)
		r.Get("/friends/requests/incoming", friendHandler.ListIncomingRequests)
		r.Post("/friends/requests/{requestID}", friendHandler.RespondToRequest)
		r.Delete("/friends/{username}", friendHandler.Unfriend)

		// Group routes
		r.Post("/groups", groupHandler.CreateGroup)
		r.Get("/groups/search", groupHandler.SearchGroups)
		r.Get("/groups/{groupID}", groupHandler.GetGroupDetails)
		r.Post("/groups/{groupID}/join", groupHandler.JoinGroup)
		r.Post("/groups/{groupID}/leave", groupHandler.LeaveGroup)

		// Message routes
		r.Post("/conversations/{conversationID}/messages", messageHandler.SendMessage)
		r.Get("/conversations/{conversationID}/messages", messageHandler.GetMessages)
		r.Post("/media/presigned-url", messageHandler.GetPresignedURL)

		// Notification routes
		r.Get("/notifications", notificationHandler.GetNotifications)
		r.Get("/notifications/unread-count", notificationHandler.GetUnreadCount)
		r.Post("/notifications/read-all", notificationHandler.MarkAllAsRead)
		r.Post("/notifications/{notificationID}/read", notificationHandler.MarkAsRead)

		// WebSocket route
		r.Get("/ws", messageHandler.ServeWs)
	})

	// Serve frontend
	fs := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fs)

	return r
}
