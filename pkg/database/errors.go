package database

import (
	"net/http"

	"github.com/srcndev/message-service/pkg/customerror"
)

// Error codes
const (
	ErrCodeDatabaseConnectionFailed = "DB_CONNECTION_FAILED"
	ErrCodeDatabaseInstanceFailed   = "DB_INSTANCE_FAILED"
	ErrCodeDatabasePingFailed       = "DB_PING_FAILED"
	ErrCodeDatabaseMigrationFailed  = "DB_MIGRATION_FAILED"
)

// Error messages
const (
	MsgDatabaseConnectionFailed = "Failed to connect to database"
	MsgDatabaseInstanceFailed   = "Failed to get database instance"
	MsgDatabasePingFailed       = "Failed to ping database"
	MsgDatabaseMigrationFailed  = "Failed to migrate database tables"
)

// Predefined errors
var (
	ErrDatabaseConnectionFailed = customerror.New(
		ErrCodeDatabaseConnectionFailed,
		MsgDatabaseConnectionFailed,
		http.StatusInternalServerError,
	)

	ErrDatabaseInstanceFailed = customerror.New(
		ErrCodeDatabaseInstanceFailed,
		MsgDatabaseInstanceFailed,
		http.StatusInternalServerError,
	)

	ErrDatabasePingFailed = customerror.New(
		ErrCodeDatabasePingFailed,
		MsgDatabasePingFailed,
		http.StatusInternalServerError,
	)

	ErrDatabaseMigrationFailed = customerror.New(
		ErrCodeDatabaseMigrationFailed,
		MsgDatabaseMigrationFailed,
		http.StatusInternalServerError,
	)
)
