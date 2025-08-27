package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Port      string
	JWTSecret string
	DBPath    string
}

// Load loads configuration from environment variables or defaults
func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "3000"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		DBPath:    getEnv("DB_PATH", "users.db"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
