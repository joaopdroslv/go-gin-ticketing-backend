package config

import "os"

type Config struct {
	HttpPort    string
	DatabaseUrl string
}

func Load() *Config {
	return &Config{
		HttpPort:    getEnv("HTTP_PORT", ":8080"),
		DatabaseUrl: getEnv("DATABASE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
