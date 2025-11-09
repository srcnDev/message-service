package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	AppURL  string

	Database DatabaseConfig
	Webhook  WebhookConfig
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

// WebhookConfig holds webhook client settings
type WebhookConfig struct {
	BaseURL    string
	AuthKey    string
	Timeout    time.Duration
	MaxRetries int
}

func NewConfig() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		fmt.Println("Info: .env file not found, using system environment variables")
	}

	webhookTimeout := 30 * time.Second
	if timeoutStr := getEnv("WEBHOOK_TIMEOUT", ""); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			webhookTimeout = timeout
		}
	}

	webhookMaxRetries := 3
	if retriesStr := getEnv("WEBHOOK_MAX_RETRIES", ""); retriesStr != "" {
		if retries, err := strconv.Atoi(retriesStr); err == nil {
			webhookMaxRetries = retries
		}
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

		Webhook: WebhookConfig{
			BaseURL:    getEnv("WEBHOOK_BASE_URL", "https://webhook.site/c3f13233-1ed4-429e-9649-8133b3b9c9cd"),
			AuthKey:    getEnv("WEBHOOK_AUTH_KEY", "INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywo"),
			Timeout:    webhookTimeout,
			MaxRetries: webhookMaxRetries,
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
	if c.Webhook.BaseURL == "" {
		return errWebhookURLEmpty
	}
	if c.Webhook.AuthKey == "" {
		return errWebhookAuthKeyEmpty
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
