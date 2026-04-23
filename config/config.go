package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port string
	Env  string

	// Database
	DatabaseURL string

	// JWT
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func Load() *Config {
	// Load .env jika ada (aman untuk production, tidak wajib)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	accessMinutes, err := strconv.Atoi(getEnv("JWT_ACCESS_MINUTES", "4320")) // 3 hari
	if err != nil {
		accessMinutes = 4320
	}

	refreshDays, err := strconv.Atoi(getEnv("JWT_REFRESH_DAYS", "30"))
	if err != nil {
		refreshDays = 30
	}

	return &Config{
		Port:                 getEnv("PORT", "8080"),
		Env:                  getEnv("APP_ENV", "development"),
		DatabaseURL:          mustEnv("DATABASE_URL"),
		JWTSecret:            mustEnv("JWT_SECRET"),
		AccessTokenDuration:  time.Duration(accessMinutes) * time.Minute,
		RefreshTokenDuration: time.Duration(refreshDays) * 24 * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return v
}
