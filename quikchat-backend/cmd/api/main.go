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
	slog.SetDefault(log)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("database connection established")

	// WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// Repositories
	userRepo := postgres.NewUserRepository(db)
	friendRepo := postgres.NewFriendRepository(db)
	blockRepo := postgres.NewBlockRepository(db)
	convoRepo := postgres.NewConversationRepository(db)
	messageRepo := postgres.NewMessageRepository(db)

	// External Services
	s3Client := s3.NewMockS3Client()

	// Use Cases / Services
	userService := service.NewUserService(userRepo, blockRepo, cfg.JWTSecretKey, cfg.AccessTokenExpiry, cfg.RefreshTokenExpiry)
	friendService := service.NewFriendService(friendRepo, userRepo, blockRepo)
	messageService := service.NewMessageService(messageRepo, convoRepo, userRepo, s3Client, hub)

	// HTTP Handlers
	userHandler := handler.NewUserHandler(userService)
	friendHandler := handler.NewFriendHandler(friendService)
	messageHandler := handler.NewMessageHandler(messageService, hub)

	// Router
	r := router.New(userHandler, friendHandler, messageHandler, hub, cfg.JWTSecretKey)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		slog.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	slog.Info("starting server", "addr", srv.Addr)

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}

	err = <-shutdownError
	if err != nil {
		slog.Error("shutdown error", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}

