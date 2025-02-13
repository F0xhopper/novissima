package config

import (
	"os"
)

type Config struct {
	WhatsAppToken string
	DatabaseURL   string
	// Add other configuration fields
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		WhatsAppToken: os.Getenv("WHATSAPP_TOKEN"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
	}, nil
} 