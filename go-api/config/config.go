package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DatabaseURL string
	JWTSecret string
	AIServiceURL string
}

func Load() *Config {
	godotenv.Load()
	return &Config {
		Port: getEnv("PORT", "8080"),
		DatabaseURl: getEnv("DATABASE_URL", ""),
		JWTSecret: 
	}
}