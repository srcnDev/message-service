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
	Redis         RedisConfig
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

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Enabled  bool
}

// WebhookConfig holds webhook client settings
type WebhookConfig struct {
	URL        string
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

	// Redis DB number
	redisDB := 0
	if dbStr := getEnv("REDIS_DB", ""); dbStr != "" {
		if db, err := strconv.Atoi(dbStr); err == nil && db >= 0 {
			redisDB = db
		}
	}

	// Redis enabled flag (default: false for optional usage)
	redisEnabled := getEnv("REDIS_ENABLED", "false") == "true"

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

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
			Enabled:  redisEnabled,
		},

		Webhook: WebhookConfig{
			URL:        getEnv("WEBHOOK_URL", "https://webhook.site/7d2fa94f-bb3c-47d7-b787-8aaacbd5097d"),
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
		return ErrAppPortEmpty
	}
	if c.AppURL == "" {
		return ErrAppURLEmpty
	}
	if c.Database.Host == "" {
		return ErrDBHostEmpty
	}
	if c.Database.Port == "" {
		return ErrDBPortEmpty
	}
	if c.Database.Username == "" {
		return ErrDBUsernameEmpty
	}
	if c.Database.Password == "" {
		return ErrDBPasswordEmpty
	}
	if c.Database.Name == "" {
		return ErrDBNameEmpty
	}
	if c.Webhook.URL == "" {
		return ErrWebhookURLEmpty
	}
	if c.Webhook.AuthKey == "" {
		return ErrWebhookAuthKeyEmpty
	}
	if c.MessageSender.Interval <= 0 {
		return ErrSenderIntervalInvalid
	}
	if c.MessageSender.BatchSize <= 0 {
		return ErrSenderBatchSizeInvalid
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
