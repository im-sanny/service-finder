package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL string
}

type AuthConfig struct {
	JWTSecret string
}

func New() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Server:   ServerConfig{Port: getEnv("PORT", "8080")},
		Database: DatabaseConfig{URL: mustGetEnv("DB_URL")},
		Auth:     AuthConfig{JWTSecret: mustGetEnv("JWT_SECRET")},
	}

	if len(cfg.Auth.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	panic(fmt.Sprintf("required env var %s is not set", key))
}
