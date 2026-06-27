package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	GinMode     string
}

// MinJWTSecretLen is the minimum required length for JWT_SECRET.
// 32 bytes = 256 bits of entropy, matching HS256's signature size.
const MinJWTSecretLen = 32

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		GinMode:     getEnv("GIN_MODE", "debug"),
	}

	if cfg.JWTSecret == "" || len(cfg.JWTSecret) < MinJWTSecretLen {
		log.Fatalf("JWT_SECRET must be set and at least %d characters (got %d)",
			MinJWTSecretLen, len(cfg.JWTSecret))
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}