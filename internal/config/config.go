package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port string
	DatabaseURL string
}

func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL belum di set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8005"
	}

	return  &Config{
		Port: port,
		DatabaseURL: dbURL}, nil
}