package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	JWTSecret   string
}

func LoadConfig() Config {
	_ = godotenv.Load()

	cfg := Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "host=db user=postgres password=postgres dbname=fleetify port=5432 sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
