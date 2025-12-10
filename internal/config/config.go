package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	ConnectionString string
}

type Config struct {
	Server ServerConfig
	DB     DatabaseConfig
}

func NewConfig(logger *slog.Logger) *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Warn("Error loading .env file", "error", err)
	}

	return &Config{
		Server: ServerConfig{
			Addr:         getEnvOrDefault("SERVER_ADDR", ":8080"),
			ReadTimeout:  getDurationEnvOrDefault("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnvOrDefault("SERVER_WRITE_TIMEOUT", 10*time.Second),
		},
		DB: DatabaseConfig{
			ConnectionString: getEnvOrDefault("GOOSE_DBSTRING", "database.db"),
		},
	}
}

func getEnvOrDefault(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		d, err := time.ParseDuration(value)
		if err != nil {
			return defaultValue
		}
		return d
	}
	return defaultValue
}
