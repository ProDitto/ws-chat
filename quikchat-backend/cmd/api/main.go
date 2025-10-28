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

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"quikchat/internal/adapter/handler/http/handler"
	"quikchat/internal/adapter/handler/http/router"
	"quikchat/internal/adapter/storage/postgres"
	"quikchat/internal/service"
	"quikchat/pkg/config"
	"quikchat/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(os.Stdout)
	slog.SetDefault(log)

	// Connect to database
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Error("unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	log.Info("successfully connected to database")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(dbpool)

	// Initialize services (use cases)
	userService := service.NewUserService(userRepo, cfg.JWTSecretKey, cfg.AccessTokenExpiry, cfg.RefreshTokenExpiry)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

	// Initialize router
	r := router.New(userHandler, cfg.JWTSecretKey)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	go func() {
		log.Info(fmt.Sprintf("starting server on port %d", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("could not start server", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Info("shutting down server gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed", "error", err)
	}
}

