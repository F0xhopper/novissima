package config

import (
	"os"

	"github.com/joho/godotenv"
)	

type Config struct {
	WhatsAppToken string
	DatabaseURL   string	
	SupabaseURL   string
	SupabaseKey   string
	SupabaseEmail string
	SupabasePassword string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		WhatsAppToken: os.Getenv("WHATSAPP_TOKEN"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		SupabaseURL:   os.Getenv("SUPABASE_URL"),
		SupabaseKey:   os.Getenv("SUPABASE_KEY"),
		SupabaseEmail: os.Getenv("SUPABASE_EMAIL"),
		SupabasePassword: os.Getenv("SUPABASE_PASSWORD"),
	}, nil
} 