package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	AppURL  string
}

func NewConfig() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		fmt.Println("Info: .env file not found, using system environment variables")
	}

	cfg := &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppURL:  getEnv("APP_URL", "http://localhost:8080"),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	if c.AppPort == "" {
		return errors.New("APP_PORT cannot be empty")
	}
	if c.AppURL == "" {
		return errors.New("APP_URL cannot be empty")
	}
	return nil
}

// getEnv reads an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
