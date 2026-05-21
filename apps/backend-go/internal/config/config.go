package config

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultHTTPAddr        = ":8080"
	defaultJWTExpiresInSec = 3600
)

type Config struct {
	HTTPAddr     string
	DatabaseURL  string
	JWTSecret    string
	JWTExpiresIn time.Duration
}

func Load() Config {
	return Config{
		HTTPAddr:     envOrDefault("HTTP_ADDR", defaultHTTPAddr),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		JWTExpiresIn: durationFromSeconds("JWT_EXPIRES_IN_SECONDS", defaultJWTExpiresInSec),
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func durationFromSeconds(key string, fallback int) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return time.Duration(fallback) * time.Second
	}
	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return time.Duration(fallback) * time.Second
	}
	return time.Duration(seconds) * time.Second
}
