package redis

import (
	"net/http"

	"github.com/srcndev/message-service/pkg/customerror"
)

// Error codes
const (
	ErrCodeRedisConnectionFailed = "REDIS_CONNECTION_FAILED"
	ErrCodeRedisPingFailed       = "REDIS_PING_FAILED"
	ErrCodeRedisSetFailed        = "REDIS_SET_FAILED"
	ErrCodeRedisGetFailed        = "REDIS_GET_FAILED"
	ErrCodeRedisDelFailed        = "REDIS_DEL_FAILED"
	ErrCodeRedisKeyNotFound      = "REDIS_KEY_NOT_FOUND"
)

// Error messages
const (
	MsgRedisConnectionFailed = "Failed to connect to Redis"
	MsgRedisPingFailed       = "Redis ping failed"
	MsgRedisSetFailed        = "Failed to set value in Redis"
	MsgRedisGetFailed        = "Failed to get value from Redis"
	MsgRedisDelFailed        = "Failed to delete key from Redis"
	MsgRedisKeyNotFound      = "Key not found in Redis"
)

// Predefined errors
var (
	ErrRedisConnectionFailed = customerror.NewCustomError(
		ErrCodeRedisConnectionFailed,
		MsgRedisConnectionFailed,
		http.StatusInternalServerError,
	)

	ErrRedisPingFailed = customerror.NewCustomError(
		ErrCodeRedisPingFailed,
		MsgRedisPingFailed,
		http.StatusInternalServerError,
	)

	ErrRedisSetFailed = customerror.NewCustomError(
		ErrCodeRedisSetFailed,
		MsgRedisSetFailed,
		http.StatusInternalServerError,
	)

	ErrRedisGetFailed = customerror.NewCustomError(
		ErrCodeRedisGetFailed,
		MsgRedisGetFailed,
		http.StatusInternalServerError,
	)

	ErrRedisDelFailed = customerror.NewCustomError(
		ErrCodeRedisDelFailed,
		MsgRedisDelFailed,
		http.StatusInternalServerError,
	)

	ErrRedisKeyNotFound = customerror.NewCustomError(
		ErrCodeRedisKeyNotFound,
		MsgRedisKeyNotFound,
		http.StatusNotFound,
	)
)
