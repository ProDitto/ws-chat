package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"quikchat/internal/adapter/handler/http/handler"
	"quikchat/internal/adapter/handler/http/router"
	"quikchat/internal/adapter/storage/postgres"
	"quikchat/internal/service"
	"quikchat/pkg/config"
	"quikchat/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log := logger.New(os.Stdout)
	log.Info("starting quikchat server")

	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Database connection
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Error("unable to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer dbpool.Close()

	// Repositories
	userRepo := postgres.NewUserRepository(dbpool)
	friendRepo := postgres.NewFriendRepository(dbpool)
	blockRepo := postgres.NewBlockRepository(dbpool)

	// Services (Use Cases)
	userService := service.NewUserService(userRepo, blockRepo, cfg.JWTSecretKey, cfg.AccessTokenExpiry, cfg.RefreshTokenExpiry)
	friendService := service.NewFriendService(friendRepo, userRepo, blockRepo)

	// HTTP Handlers
	userHandler := handler.NewUserHandler(userService)
	friendHandler := handler.NewFriendHandler(friendService)

	// Router
	r := router.New(userHandler, friendHandler, cfg.JWTSecretKey)

	// Server setup
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("server exited properly")
}

