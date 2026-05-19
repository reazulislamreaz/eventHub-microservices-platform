package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPPort          string
	SwaggerHost       string // optional; empty = same host as browser (recommended)
	JWTSecret         string
	JWTExpiration     time.Duration
	UserServiceAddr   string
	EventServiceAddr  string
	TicketServiceAddr string
}

func Load() *Config {
	expHours := 24
	return &Config{
		HTTPPort:          getEnv("HTTP_PORT", "8080"),
		SwaggerHost:       getEnv("SWAGGER_HOST", ""),
		JWTSecret:         getEnv("JWT_SECRET", "eventhub-dev-secret-change-in-production"),
		JWTExpiration:     time.Duration(expHours) * time.Hour,
		UserServiceAddr:   getEnv("USER_SERVICE_ADDR", "localhost:50051"),
		EventServiceAddr:  getEnv("EVENT_SERVICE_ADDR", "localhost:50052"),
		TicketServiceAddr: getEnv("TICKET_SERVICE_ADDR", "localhost:50053"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
