package config

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type AppConfig struct {
	DB    *sql.DB
	DbUrl string
}

func LoadConfig() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dbUrl := getEnv("DB_URL", "postgres://postgres:password@localhost:5432/database?sslmode=disable")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &AppConfig{
		DB:    db,
		DbUrl: dbUrl,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
