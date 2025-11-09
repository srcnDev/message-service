package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	AppURL  string

	Database DatabaseConfig
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

func NewConfig() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		fmt.Println("Info: .env file not found, using system environment variables")
	}

	cfg := &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppURL:  getEnv("APP_URL", "http://localhost:8080"),

		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_DB_HOST", "localhost"),
			Port:     getEnv("POSTGRES_DB_PORT", "5432"),
			Username: getEnv("POSTGRES_DB_USERNAME", "user"),
			Password: getEnv("POSTGRES_DB_PASSWORD", "password"),
			Name:     getEnv("POSTGRES_DB_NAME", "messagedb"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.AppPort == "" {
		return errAppPortEmpty
	}
	if c.AppURL == "" {
		return errAppURLEmpty
	}
	if c.Database.Host == "" {
		return errDBHostEmpty
	}
	if c.Database.Port == "" {
		return errDBPortEmpty
	}
	if c.Database.Username == "" {
		return errDBUsernameEmpty
	}
	if c.Database.Password == "" {
		return errDBPasswordEmpty
	}
	if c.Database.Name == "" {
		return errDBNameEmpty
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
