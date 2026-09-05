package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	HTTPAddr    string
	CORSOrigins []string
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		DatabaseURL: strings.TrimSpace(os.Getenv("DATABASE_URL")),
		HTTPAddr:    strings.TrimSpace(os.Getenv("HTTP_ADDR")),
		CORSOrigins: defaultCORSOrigins(),
	}
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "postgres://kc:kc@127.0.0.1:5434/knowledge_center?sslmode=disable"
	}
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":8091"
	}
	if extra := splitCSV(os.Getenv("CORS_ORIGINS")); len(extra) > 0 {
		cfg.CORSOrigins = uniqueStrings(append(cfg.CORSOrigins, extra...))
	}
	return cfg
}

func defaultCORSOrigins() []string {
	return []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:4173",
		"http://127.0.0.1:4173",
		"http://localhost:4006",
		"http://127.0.0.1:4006",
	}
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func uniqueStrings(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}
