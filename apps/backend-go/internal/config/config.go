package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHTTPAddr        = ":8080"
	defaultJWTExpiresInSec = 3600
)

type Config struct {
	HTTPAddr       string
	DatabaseURL    string
	JWTSecret      string
	JWTExpiresIn   time.Duration
	AllowedOrigins []string
	AIProvider     string
	AIModel        string
	OpenAIKey      string
	GeminiKey      string
	DeepSeekKey    string
}

func (c Config) Validate() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}
	if len(c.JWTSecret) < 32 {
		return errors.New("JWT_SECRET must be at least 32 characters long for security")
	}
	return nil
}

func Load() Config {
	return Config{
		HTTPAddr:       envOrDefault("HTTP_ADDR", defaultHTTPAddr),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpiresIn:   durationFromSeconds("JWT_EXPIRES_IN_SECONDS", defaultJWTExpiresInSec),
		AllowedOrigins: strings.Split(envOrDefault("ALLOWED_ORIGINS", "*"), ","),
		AIProvider:     envOrDefault("AI_PROVIDER", "openai"),
		AIModel:        os.Getenv("AI_MODEL"),
		OpenAIKey:      os.Getenv("OPENAI_API_KEY"),
		GeminiKey:      os.Getenv("GEMINI_API_KEY"),
		DeepSeekKey:    os.Getenv("DEEPSEEK_API_KEY"),
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
