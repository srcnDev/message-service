package config

import (
	"github.com/srcndev/message-service/pkg/customerror"
)

var (
	errAppPortEmpty = customerror.NewWithDefaults("APP_PORT_EMPTY", "APP_PORT cannot be empty")
	errAppURLEmpty  = customerror.NewWithDefaults("APP_URL_EMPTY", "APP_URL cannot be empty")

	errDBHostEmpty     = customerror.NewWithDefaults("DB_HOST_EMPTY", "Database host cannot be empty")
	errDBPortEmpty     = customerror.NewWithDefaults("DB_PORT_EMPTY", "Database port cannot be empty")
	errDBUsernameEmpty = customerror.NewWithDefaults("DB_USERNAME_EMPTY", "Database username cannot be empty")
	errDBPasswordEmpty = customerror.NewWithDefaults("DB_PASSWORD_EMPTY", "Database password cannot be empty")
	errDBNameEmpty     = customerror.NewWithDefaults("DB_NAME_EMPTY", "Database name cannot be empty")

	errWebhookURLEmpty     = customerror.NewWithDefaults("WEBHOOK_URL_EMPTY", "Webhook URL cannot be empty")
	errWebhookAuthKeyEmpty = customerror.NewWithDefaults("WEBHOOK_AUTH_KEY_EMPTY", "Webhook auth key cannot be empty")
)
