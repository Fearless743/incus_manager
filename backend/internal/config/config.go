package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	IncusURL    string
	IncusCert   string
	JWTSecret   string
	Port        string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost/incus_manager"),
		IncusURL:    getEnv("INCUS_URL", "https://localhost:8443"),
		IncusCert:   getEnv("INCUS_CERT", ""),
		JWTSecret:   getEnv("JWT_SECRET", "change-this-secret-in-production"),
		Port:        getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
