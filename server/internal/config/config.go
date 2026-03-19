package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env  string
	Port string

	Firebase FirebaseConfig
}

type FirebaseConfig struct {
	ProjectID       string
	CredentialsFile string
}

func Load() (*Config, error) {

	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Env:  getEnv("ENV", "development"),
		Port: getEnv("APP_PORT", "3000"),

		Firebase: FirebaseConfig{
			ProjectID:       getEnv("FIREBASE_PROJECT_ID", ""),
			CredentialsFile: getEnv("FIREBASE_CREDENTIALS_FILE", ""),
		},
	}

	if err := cfg.Firebase.Normalize(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (f *FirebaseConfig) Normalize() error {
	if f.ProjectID == "" {
		return fmt.Errorf("FIREBASE_PROJECT_ID is required")
	}
	if f.CredentialsFile == "" {
		return fmt.Errorf("FIREBASE_CREDENTIALS_FILE is required")
	}
	return nil
}

func getEnv(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
