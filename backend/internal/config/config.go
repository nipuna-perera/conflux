// Configuration management for the backend application
// Loads and validates environment variables and application settings
// Provides centralized configuration access throughout the application
package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	// Server configuration
	Port string
	Host string

	// Database configuration
	DBType     string // "mysql" or "postgres"
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string

	// JWT configuration
	JWTSecret     string
	JWTExpiration int

	// CORS configuration
	AllowedOrigins []string
}

// Load reads configuration from environment variables
// Validates required settings and returns configured struct
func Load() (*Config, error) {
	config := &Config{
		Port:       getEnv("PORT", "8080"),
		Host:       getEnv("HOST", "0.0.0.0"),
		DBType:     getEnv("DB_TYPE", "mysql"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "appdb"),
		DBUser:     getEnv("DB_USER", "appuser"),
		DBPassword: getEnv("DB_PASSWORD", "apppassword"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
	}

	// Parse JWT expiration
	expStr := getEnv("JWT_EXPIRATION", "3600")
	if exp, err := strconv.Atoi(expStr); err == nil {
		config.JWTExpiration = exp
	} else {
		config.JWTExpiration = 3600
	}

	// Parse allowed origins
	originsStr := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")
	config.AllowedOrigins = strings.Split(originsStr, ",")

	return config, nil
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
