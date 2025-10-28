package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"quikchat/internal/adapter/external/s3"
	"quikchat/internal/adapter/handler/http/handler"
	"quikchat/internal/adapter/handler/http/router"
	"quikchat/internal/adapter/storage/postgres"
	"quikchat/internal/adapter/ws"
	"quikchat/internal/service"
	"quikchat/pkg/config"
	"quikchat/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Initialize logger
	log := logger.New(os.Stdout)
	slog.SetDefault(log)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Connect to database
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	slog.Info("successfully connected to database")

	// Initialize WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(dbpool)
	friendRepo := postgres.NewFriendRepository(dbpool)
	blockRepo := postgres.NewBlockRepository(dbpool)
	convoRepo := postgres.NewConversationRepository(dbpool)
	msgRepo := postgres.NewMessageRepository(dbpool)
	groupRepo := postgres.NewGroupRepository(dbpool)
	notificationRepo := postgres.NewNotificationRepository(dbpool)

	// Initialize external services
	s3Client := s3.NewMockS3Client() // Replace with real S3 client in production

	// Initialize services (use cases)
	notificationService := service.NewNotificationService(notificationRepo, userRepo, hub)
	userService := service.NewUserService(userRepo, blockRepo, cfg.JWTSecretKey, cfg.AccessTokenExpiry, cfg.RefreshTokenExpiry)
	friendService := service.NewFriendService(friendRepo, userRepo, blockRepo, notificationService)
	messageService := service.NewMessageService(msgRepo, convoRepo, userRepo, notificationService, s3Client, hub)
	groupService := service.NewGroupService(groupRepo, userRepo, convoRepo, notificationService, cfg.MaxGroupMembers)

	// Initialize HTTP handlers
	userHandler := handler.NewUserHandler(userService)
	friendHandler := handler.NewFriendHandler(friendService)
	messageHandler := handler.NewMessageHandler(messageService, hub)
	groupHandler := handler.NewGroupHandler(groupService)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	// Initialize router
	r := router.New(userHandler, friendHandler, messageHandler, groupHandler, notificationHandler, hub, cfg.JWTSecretKey)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	go func() {
		slog.Info(fmt.Sprintf("starting server on port %d", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited properly")
}

