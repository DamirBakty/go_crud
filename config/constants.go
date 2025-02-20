package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbUrl string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DbUrl: getEnv("DB_URL", "postgres://postgres:password@localhost:5432/database?sslmode=disable"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
