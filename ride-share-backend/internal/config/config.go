package config

import (
	"os"
)

type Config struct {
	DatabaseURL         string
	JWTSecret           string
	MpesaConsumerKey    string
	MpesaConsumerSecret string
	MpesaPasskey        string
	MpesaShortcode      string
	MpesaCallbackURL    string
	Environment         string
	// URL Configuration
	BaseURL         string
	APIBasePath     string
	// Note: No frontend URL needed - Flutter mobile app communicates directly with API
}

func Load() *Config {
	cfg := &Config{
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://user:password@localhost/kenyan_ride_share?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		MpesaConsumerKey:   getEnv("MPESA_CONSUMER_KEY", ""),
		MpesaConsumerSecret: getEnv("MPESA_CONSUMER_SECRET", ""),
		MpesaPasskey:       getEnv("MPESA_PASSKEY", ""),
		MpesaShortcode:     getEnv("MPESA_SHORTCODE", ""),
		Environment:        getEnv("ENVIRONMENT", "development"),
		// URL Configuration
		BaseURL:         getEnv("BASE_URL", "http://localhost:8080"),
		APIBasePath:     getEnv("API_BASE_PATH", "/api/v1"),
	}
	
	// Build M-Pesa callback URL dynamically if not explicitly set
	cfg.MpesaCallbackURL = getEnv("MPESA_CALLBACK_URL", cfg.GetMpesaCallbackURL())
	
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetAPIURL returns the full API URL with base path
func (c *Config) GetAPIURL() string {
	return c.BaseURL + c.APIBasePath
}

// GetCallbackURL returns the full M-Pesa callback URL
func (c *Config) GetMpesaCallbackURL() string {
	return c.BaseURL + c.APIBasePath + "/payments/mpesa/callback"
}

