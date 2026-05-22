package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	AIServiceURL string
}

func Load() *Config {
	godotenv.Load()
	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		JWTSecret:    getEnv("JWT_SECRET", ""),
		AIServiceURL: getEnv("AI_SERVICE_URL", "http://localhost:8000"),
	}
}

// parameters: 'key' will be our key from a mapping in the env file and 'fallback' is there incase the key has no value
func getEnv(key, fallback string) string {
	/*
		Go uses this format for if statements:

		if INITIALIZATION; CONDITION {
			LOGIC
		}

		The os.Getenv() returns the value for the respective key we provide so in this case
		our code will store the value of the env key into value 'val'

	*/
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
