package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client interface defines Redis operations
type Client interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
	Close() error
	Ping(ctx context.Context) error
}

// Config holds Redis connection settings
type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// client is the private implementation of Client interface
type client struct {
	rdb *redis.Client
}

// Compile-time interface compliance check
var _ Client = (*client)(nil)
