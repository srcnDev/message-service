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
	ErrAppPortEmpty = customerror.New(
		ErrCodeAppPortEmpty,
		MsgAppPortEmpty,
		http.StatusBadRequest,
	)

	ErrAppURLEmpty = customerror.New(
		ErrCodeAppURLEmpty,
		MsgAppURLEmpty,
		http.StatusBadRequest,
	)

	ErrDBHostEmpty = customerror.New(
		ErrCodeDBHostEmpty,
		MsgDBHostEmpty,
		http.StatusBadRequest,
	)

	ErrDBPortEmpty = customerror.New(
		ErrCodeDBPortEmpty,
		MsgDBPortEmpty,
		http.StatusBadRequest,
	)

	ErrDBUsernameEmpty = customerror.New(
		ErrCodeDBUsernameEmpty,
		MsgDBUsernameEmpty,
		http.StatusBadRequest,
	)

	ErrDBPasswordEmpty = customerror.New(
		ErrCodeDBPasswordEmpty,
		MsgDBPasswordEmpty,
		http.StatusBadRequest,
	)

	ErrDBNameEmpty = customerror.New(
		ErrCodeDBNameEmpty,
		MsgDBNameEmpty,
		http.StatusBadRequest,
	)

	ErrWebhookURLEmpty = customerror.New(
		ErrCodeWebhookURLEmpty,
		MsgWebhookURLEmpty,
		http.StatusBadRequest,
	)

	ErrWebhookAuthKeyEmpty = customerror.New(
		ErrCodeWebhookAuthKeyEmpty,
		MsgWebhookAuthKeyEmpty,
		http.StatusBadRequest,
	)

	ErrSenderIntervalInvalid = customerror.New(
		ErrCodeSenderIntervalInvalid,
		MsgSenderIntervalInvalid,
		http.StatusBadRequest,
	)

	ErrSenderBatchSizeInvalid = customerror.New(
		ErrCodeSenderBatchSizeInvalid,
		MsgSenderBatchSizeInvalid,
		http.StatusBadRequest,
	)
)
