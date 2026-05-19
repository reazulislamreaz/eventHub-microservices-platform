package config

import (
	"fmt"
	"os"
)

type Config struct {
	GRPCPort       string
	EventServiceAddr string
	DB             DBConfig
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
	return &Config{
		GRPCPort:         getEnv("GRPC_PORT", "50053"),
		EventServiceAddr: getEnv("EVENT_SERVICE_ADDR", "localhost:50052"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5434"),
			User:     getEnv("DB_USER", "eventhub"),
			Password: getEnv("DB_PASSWORD", "eventhub_secret"),
			Name:     getEnv("DB_NAME", "ticket_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}, nil
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
