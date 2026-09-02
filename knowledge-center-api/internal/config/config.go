package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	HTTPAddr    string
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		DatabaseURL: strings.TrimSpace(os.Getenv("DATABASE_URL")),
		HTTPAddr:    strings.TrimSpace(os.Getenv("HTTP_ADDR")),
	}
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "postgres://kc:kc@127.0.0.1:5434/knowledge_center?sslmode=disable"
	}
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":8091"
	}
	return cfg
}
