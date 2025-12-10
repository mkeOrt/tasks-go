package config

import (
	"os"
	"time"
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

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Addr:         ":8080",
			ReadTimeout:  GetDurationEnvOrDefault("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: GetDurationEnvOrDefault("SERVER_WRITE_TIMEOUT", 10*time.Second),
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

func GetDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		d, err := time.ParseDuration(value)
		if err != nil {
			return defaultValue
		}
		return d
	}
	return defaultValue
}
