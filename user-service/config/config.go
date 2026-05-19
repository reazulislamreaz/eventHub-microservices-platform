package config

import (
	"fmt"
	"os"
)

type Config struct {
	GRPCPort string
	DB       DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() (*Config, error) {
	cfg := &Config{
		GRPCPort: getEnv("GRPC_PORT", "50051"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "eventhub"),
			Password: getEnv("DB_PASSWORD", "eventhub_secret"),
			Name:     getEnv("DB_NAME", "user_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
	return cfg, nil
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
