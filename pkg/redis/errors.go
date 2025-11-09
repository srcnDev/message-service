package redis

import "github.com/srcndev/message-service/pkg/customerror"

var (
	// ErrRedisConnectionFailed indicates Redis connection failed
	ErrRedisConnectionFailed = customerror.New("REDIS_CONNECTION_FAILED", "Failed to connect to Redis", 500)

	// ErrRedisPingFailed indicates Redis ping failed
	ErrRedisPingFailed = customerror.New("REDIS_PING_FAILED", "Redis ping failed", 500)

	// ErrRedisSetFailed indicates Redis set operation failed
	ErrRedisSetFailed = customerror.New("REDIS_SET_FAILED", "Failed to set value in Redis", 500)

	// ErrRedisGetFailed indicates Redis get operation failed
	ErrRedisGetFailed = customerror.New("REDIS_GET_FAILED", "Failed to get value from Redis", 500)

	// ErrRedisDelFailed indicates Redis delete operation failed
	ErrRedisDelFailed = customerror.New("REDIS_DEL_FAILED", "Failed to delete key from Redis", 500)

	// ErrRedisKeyNotFound indicates the key was not found in Redis
	ErrRedisKeyNotFound = customerror.New("REDIS_KEY_NOT_FOUND", "Key not found in Redis", 404)
)
