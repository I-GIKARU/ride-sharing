package config

import (
	"os"
)

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	MpesaConsumerKey    string
	MpesaConsumerSecret string
	MpesaPasskey        string
	MpesaShortcode      string
	MpesaCallbackURL    string
	Environment         string
}

func Load() *Config {
	return &Config{
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://user:password@localhost/kenyan_ride_share?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		MpesaConsumerKey:   getEnv("MPESA_CONSUMER_KEY", ""),
		MpesaConsumerSecret: getEnv("MPESA_CONSUMER_SECRET", ""),
		MpesaPasskey:       getEnv("MPESA_PASSKEY", ""),
		MpesaShortcode:     getEnv("MPESA_SHORTCODE", ""),
		MpesaCallbackURL:   getEnv("MPESA_CALLBACK_URL", ""),
		Environment:        getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

