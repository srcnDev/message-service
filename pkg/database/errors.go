package database

import "github.com/srcndev/message-service/pkg/customerror"

var (
	errDatabaseConnectionFailed = customerror.NewWithDefaults("DB_CONNECTION_FAILED", "Failed to connect to database")
	errDatabaseInstanceFailed   = customerror.NewWithDefaults("DB_INSTANCE_FAILED", "Failed to get database instance")
	errDatabasePingFailed       = customerror.NewWithDefaults("DB_PING_FAILED", "Failed to ping database")
)
