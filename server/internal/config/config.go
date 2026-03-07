package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env  string
	Port string

	Firebase FirebaseConfig
}

type FirebaseConfig struct {
	ProjectID    string
	DatabaseURL  string
	EmulatorHost string
}

func Load() (*Config, error) {

	// load .env if present
	err := godotenv.Load("../.env")

	cfg := &Config{
		Env:  getEnv("ENV", "development"),
		Port: getEnv("APP_PORT", "3000"),

		Firebase: FirebaseConfig{
			ProjectID:    getEnv("FIREBASE_PROJECT_ID", ""),
			DatabaseURL:  getEnv("FIREBASE_DATABASE_URL", ""),
			EmulatorHost: getEnv("FIREBASE_DATABASE_EMULATOR_HOST", ""),
		},
	}

	return cfg, err
}

func getEnv(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
