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
	log := logger.New(os.Stdout)
	slog.SetDefault(log) // Keep original slog default setting

	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load config", slog.Any("error", err)) // Adopt attempted slog.Any
		os.Exit(1)
	}

	// Connect to database
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL) // Adopt attempted variable name 'pool'
	if err != nil {
		log.Error("failed to connect to database", slog.Any("error", err)) // Adopt attempted slog.Any
		os.Exit(1)
	}
	defer pool.Close()

	log.Info("database connection established") // Adopt attempted log.Info

	// Initialize WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(pool)
	friendRepo := postgres.NewFriendRepository(pool)
	blockRepo := postgres.NewBlockRepository(pool)
	convoRepo := postgres.NewConversationRepository(pool)
	messageRepo := postgres.NewMessageRepository(pool)
	groupRepo := postgres.NewGroupRepository(pool) // Add new groupRepo

	// Initialize external services
	s3Client := s3.NewMockS3Client() // Replace with real S3 client

	// Initialize services (use cases)
	userService := service.NewUserService(userRepo, blockRepo, cfg.JWTSecretKey, cfg.AccessTokenExpiry, cfg.RefreshTokenExpiry)
	friendService := service.NewFriendService(friendRepo, userRepo, blockRepo)
	messageService := service.NewMessageService(messageRepo, convoRepo, userRepo, s3Client, hub)
	groupService := service.NewGroupService(groupRepo, userRepo, convoRepo) // Add new groupService

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	friendHandler := handler.NewFriendHandler(friendService)
	messageHandler := handler.NewMessageHandler(messageService, hub)
	groupHandler := handler.NewGroupHandler(groupService) // Add new groupHandler

	// Initialize router
	r := router.New(userHandler, friendHandler, messageHandler, groupHandler, hub, cfg.JWTSecretKey) // Add groupHandler to router.New

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
		// Original timeouts removed as per attempted content
	}

	go func() {
		log.Info("starting server", slog.Int("port", cfg.Port)) // Adopt attempted log format
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", slog.Any("error", err)) // Adopt attempted log format
			os.Exit(1)
		}
	}()

	// Graceful shutdown (Adopt attempted shutdown logic)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("server exited properly")
}

