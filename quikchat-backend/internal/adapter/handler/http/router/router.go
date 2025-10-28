package router

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

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

	// Middleware
	r.Use(chi_middleware.RequestID)
	r.Use(chi_middleware.RealIP)
	r.Use(middleware.NewStructuredLogger(slog.Default()))
	r.Use(chi_middleware.Recoverer)

	// Public routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.RateLimiter(2, 5)) // 2 req/sec, burst of 5
			r.Post("/auth/register", userHandler.Register)
			r.Post("/auth/login", userHandler.Login)
		})

		// WebSocket
		r.Get("/ws", messageHandler.ServeWs)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtSecret))

			// Users
			r.Get("/users/me", userHandler.GetCurrentUser)
			r.Get("/users/profile/{username}", userHandler.GetProfileByUsername)
			r.Group(func(r chi.Router) {
				r.Use(middleware.RateLimiter(5, 10)) // 5 req/sec, burst of 10 for profile updates
				r.Put("/users/me/profile", userHandler.UpdateProfile)
			})
			r.Post("/users/block/{username}", userHandler.BlockUser)
			r.Delete("/users/unblock/{username}", userHandler.UnblockUser)

			// Friends
			r.Route("/friends", func(r chi.Router) {
				r.Get("/", friendHandler.ListFriends)
				r.Post("/requests", friendHandler.SendRequest)
				r.Get("/requests", friendHandler.ListIncomingRequests)
				r.Post("/requests/{requestID}", friendHandler.RespondToRequest)
				r.Delete("/{username}", friendHandler.Unfriend)
			})

			// Groups
			r.Route("/groups", func(r chi.Router) {
				r.Post("/", groupHandler.CreateGroup)
				r.Get("/search", groupHandler.SearchGroups)
				r.Get("/{groupID}", groupHandler.GetGroupDetails)
				r.Post("/{groupID}/join", groupHandler.JoinGroup)
				r.Post("/{groupID}/leave", groupHandler.LeaveGroup)
			})

			// Messages
			r.Get("/conversations/{conversationID}/messages", messageHandler.GetMessages)
			r.Post("/conversations/{conversationID}/messages", messageHandler.SendMessage)
			r.Get("/media/upload-url", messageHandler.GetPresignedURL)

			// Notifications
			r.Route("/notifications", func(r chi.Router) {
				r.Get("/", notificationHandler.GetNotifications)
				r.Get("/unread-count", notificationHandler.GetUnreadCount)
				r.Post("/read-all", notificationHandler.MarkAllAsRead)
				r.Post("/{notificationID}/read", notificationHandler.MarkAsRead)
			})
		})
	})

	// Serve frontend static files
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))
	FileServer(r, "/", filesDir)

	return r
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem with SPA fallback.
func FileServer(r chi.Router, path string, root http.Dir) {
	fsRoot := string(root)

	// Check directory existence
	if _, err := os.Stat(fsRoot); os.IsNotExist(err) {
		slog.Warn("static file directory does not exist, skipping file server", "path", fsRoot)
		return
	}

	// Ensure trailing slash in route path
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}

	fileServer := http.StripPrefix(path, http.FileServer(root))

	r.Get(path+"*", func(w http.ResponseWriter, req *http.Request) {
		requested := req.URL.Path[len(path)-1:] // remove prefix
		// If file not found → serve index.html for SPA
		if _, err := os.Stat(filepath.Join(fsRoot, requested)); os.IsNotExist(err) {
			http.ServeFile(w, req, filepath.Join(fsRoot, "index.html"))
			return
		}
		fileServer.ServeHTTP(w, req)
	})
}
