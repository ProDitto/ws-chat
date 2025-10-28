package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	Port               int
	DatabaseURL        string
	RedisAddress       string
	JWTSecretKey       string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	MaxGroupMembers    int
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	_ = godotenv.Load() // Ignore error if .env file doesn't exist

	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	redisAddress := getEnv("REDIS_ADDRESS", "localhost:6379")

	jwtSecretKey := getEnv("JWT_SECRET_KEY", "")
	if jwtSecretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY is not set")
	}

	accessExpStr := getEnv("ACCESS_TOKEN_EXPIRY", "15m")
	accessExp, err := time.ParseDuration(accessExpStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ACCESS_TOKEN_EXPIRY: %w", err)
	}

	refreshExpStr := getEnv("REFRESH_TOKEN_EXPIRY", "72h")
	refreshExp, err := time.ParseDuration(refreshExpStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_EXPIRY: %w", err)
	}

	maxGroupMembersStr := getEnv("MAX_GROUP_MEMBERS", "50")
	maxGroupMembers, err := strconv.Atoi(maxGroupMembersStr)
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_GROUP_MEMBERS: %w", err)
	}

	cfg := &Config{
		Port:               port,
		DatabaseURL:        databaseURL,
		RedisAddress:       redisAddress,
		JWTSecretKey:       jwtSecretKey,
		AccessTokenExpiry:  accessExp,
		RefreshTokenExpiry: refreshExp,
		MaxGroupMembers:    maxGroupMembers,
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

