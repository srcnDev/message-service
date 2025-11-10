package config

import (
	"net/http"

	"github.com/srcndev/message-service/pkg/customerror"
)

// Error codes
const (
	ErrCodeAppPortEmpty           = "APP_PORT_EMPTY"
	ErrCodeAppURLEmpty            = "APP_URL_EMPTY"
	ErrCodeDBHostEmpty            = "DB_HOST_EMPTY"
	ErrCodeDBPortEmpty            = "DB_PORT_EMPTY"
	ErrCodeDBUsernameEmpty        = "DB_USERNAME_EMPTY"
	ErrCodeDBPasswordEmpty        = "DB_PASSWORD_EMPTY"
	ErrCodeDBNameEmpty            = "DB_NAME_EMPTY"
	ErrCodeWebhookURLEmpty        = "WEBHOOK_URL_EMPTY"
	ErrCodeWebhookAuthKeyEmpty    = "WEBHOOK_AUTH_KEY_EMPTY"
	ErrCodeSenderIntervalInvalid  = "SENDER_INTERVAL_INVALID"
	ErrCodeSenderBatchSizeInvalid = "SENDER_BATCH_SIZE_INVALID"
)

// Error messages
const (
	MsgAppPortEmpty           = "APP_PORT cannot be empty"
	MsgAppURLEmpty            = "APP_URL cannot be empty"
	MsgDBHostEmpty            = "Database host cannot be empty"
	MsgDBPortEmpty            = "Database port cannot be empty"
	MsgDBUsernameEmpty        = "Database username cannot be empty"
	MsgDBPasswordEmpty        = "Database password cannot be empty"
	MsgDBNameEmpty            = "Database name cannot be empty"
	MsgWebhookURLEmpty        = "Webhook URL cannot be empty"
	MsgWebhookAuthKeyEmpty    = "Webhook auth key cannot be empty"
	MsgSenderIntervalInvalid  = "Message sender interval must be greater than 0"
	MsgSenderBatchSizeInvalid = "Message sender batch size must be greater than 0"
)

// Predefined errors
var (
	ErrAppPortEmpty = customerror.NewCustomError(
		ErrCodeAppPortEmpty,
		MsgAppPortEmpty,
		http.StatusBadRequest,
	)

	ErrAppURLEmpty = customerror.NewCustomError(
		ErrCodeAppURLEmpty,
		MsgAppURLEmpty,
		http.StatusBadRequest,
	)

	ErrDBHostEmpty = customerror.NewCustomError(
		ErrCodeDBHostEmpty,
		MsgDBHostEmpty,
		http.StatusBadRequest,
	)

	ErrDBPortEmpty = customerror.NewCustomError(
		ErrCodeDBPortEmpty,
		MsgDBPortEmpty,
		http.StatusBadRequest,
	)

	ErrDBUsernameEmpty = customerror.NewCustomError(
		ErrCodeDBUsernameEmpty,
		MsgDBUsernameEmpty,
		http.StatusBadRequest,
	)

	ErrDBPasswordEmpty = customerror.NewCustomError(
		ErrCodeDBPasswordEmpty,
		MsgDBPasswordEmpty,
		http.StatusBadRequest,
	)

	ErrDBNameEmpty = customerror.NewCustomError(
		ErrCodeDBNameEmpty,
		MsgDBNameEmpty,
		http.StatusBadRequest,
	)

	ErrWebhookURLEmpty = customerror.NewCustomError(
		ErrCodeWebhookURLEmpty,
		MsgWebhookURLEmpty,
		http.StatusBadRequest,
	)

	ErrWebhookAuthKeyEmpty = customerror.NewCustomError(
		ErrCodeWebhookAuthKeyEmpty,
		MsgWebhookAuthKeyEmpty,
		http.StatusBadRequest,
	)

	ErrSenderIntervalInvalid = customerror.NewCustomError(
		ErrCodeSenderIntervalInvalid,
		MsgSenderIntervalInvalid,
		http.StatusBadRequest,
	)

	ErrSenderBatchSizeInvalid = customerror.NewCustomError(
		ErrCodeSenderBatchSizeInvalid,
		MsgSenderBatchSizeInvalid,
		http.StatusBadRequest,
	)
)
