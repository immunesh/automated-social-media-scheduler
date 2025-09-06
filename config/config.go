package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	DatabasePath   string
	JWTSecret     string
	Port          string
	GinMode       string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		DatabasePath: getEnv("DATABASE_PATH", "./scheduler.db"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Port:         getEnv("PORT", "8080"),
		GinMode:      getEnv("GIN_MODE", "debug"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}