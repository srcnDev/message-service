package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewClient creates a new Redis client
func NewClient(cfg Config) (Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRedisPingFailed, err)
	}

	return &client{rdb: rdb}, nil
}

// Set stores a value in Redis with expiration
func (c *client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := c.rdb.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("%w: %v", ErrRedisSetFailed, err)
	}
	return nil
}

// Get retrieves a value from Redis
func (c *client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrRedisKeyNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrRedisGetFailed, err)
	}
	return val, nil
}

// Del deletes keys from Redis
func (c *client) Del(ctx context.Context, keys ...string) error {
	if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("%w: %v", ErrRedisDelFailed, err)
	}
	return nil
}

// Exists checks if keys exist in Redis
func (c *client) Exists(ctx context.Context, keys ...string) (int64, error) {
	count, err := c.rdb.Exists(ctx, keys...).Result()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrRedisGetFailed, err)
	}
	return count, nil
}

// Close closes the Redis connection
func (c *client) Close() error {
	if err := c.rdb.Close(); err != nil {
		return fmt.Errorf("%w: %v", ErrRedisConnectionFailed, err)
	}
	return nil
}

// Ping tests the Redis connection
func (c *client) Ping(ctx context.Context) error {
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("%w: %v", ErrRedisPingFailed, err)
	}
	return nil
}
