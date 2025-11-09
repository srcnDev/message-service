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

	Database      DatabaseConfig
	Webhook       WebhookConfig
	MessageSender MessageSenderConfig
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

// MessageSenderConfig holds message sender job settings
type MessageSenderConfig struct {
	Interval  time.Duration // How often to check for pending messages
	BatchSize int           // Number of messages to send per cycle
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

	// Message sender interval (default: 120 seconds = 2 minutes as per case study)
	senderInterval := 120 * time.Second
	if intervalStr := getEnv("MESSAGE_SENDER_INTERVAL", ""); intervalStr != "" {
		if intervalSec, err := strconv.Atoi(intervalStr); err == nil && intervalSec > 0 {
			senderInterval = time.Duration(intervalSec) * time.Second
		}
	}

	// Message sender batch size (default: 2 messages per cycle as per case study)
	senderBatchSize := 2
	if batchStr := getEnv("MESSAGE_SENDER_BATCH_SIZE", ""); batchStr != "" {
		if batch, err := strconv.Atoi(batchStr); err == nil && batch > 0 {
			senderBatchSize = batch
		}
	}

	cfg := &Config{
		AppPort: getEnv("APP_PORT", "8000"),
		AppURL:  getEnv("APP_URL", "http://localhost:8000"),

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

		MessageSender: MessageSenderConfig{
			Interval:  senderInterval,
			BatchSize: senderBatchSize,
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
	if c.MessageSender.Interval <= 0 {
		return errSenderIntervalInvalid
	}
	if c.MessageSender.BatchSize <= 0 {
		return errSenderBatchSizeInvalid
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
