package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               int
	DatabaseURL        string
	RedisAddress       string
	JWTSecretKey       string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, err
	}

	accessExp, err := strconv.Atoi(getEnv("JWT_ACCESS_TOKEN_EXPIRY_MINUTES", "15"))
	if err != nil {
		return nil, err
	}

	refreshExp, err := strconv.Atoi(getEnv("JWT_REFRESH_TOKEN_EXPIRY_DAYS", "7"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:               port,
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		RedisAddress:       getEnv("REDIS_ADDRESS", ""),
		JWTSecretKey:       getEnv("JWT_SECRET_KEY", "default-secret"),
		AccessTokenExpiry:  time.Duration(accessExp) * time.Minute,
		RefreshTokenExpiry: time.Duration(refreshExp) * 24 * time.Hour,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

