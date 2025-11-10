package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/srcndev/message-service/pkg/redis"
)

// MessageCacheRepository interface defines cache operations for messages
type MessageCacheRepository interface {
	CacheSentMessage(ctx context.Context, messageID string, sentAt time.Time) error
	GetCachedMessage(ctx context.Context, messageID string) (*CachedMessage, error)
	IsCached(ctx context.Context, messageID string) (bool, error)
}

// CachedMessage represents a cached message in Redis
type CachedMessage struct {
	MessageID string    `json:"messageId"`
	SentAt    time.Time `json:"sentAt"`
}

// messageCacheRepository is the private implementation
type messageCacheRepository struct {
	redis redis.Client
}

// Compile-time interface compliance check
var _ MessageCacheRepository = (*messageCacheRepository)(nil)

// NewMessageCacheRepository creates a new message cache repository
func NewMessageCacheRepository(redisClient redis.Client) MessageCacheRepository {
	return &messageCacheRepository{
		redis: redisClient,
	}
}

// CacheSentMessage stores message send information in Redis
// Key format: message:{messageId}
// TTL: 30 days (can be adjusted)
func (r *messageCacheRepository) CacheSentMessage(ctx context.Context, messageID string, sentAt time.Time) error {
	cached := CachedMessage{
		MessageID: messageID,
		SentAt:    sentAt,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return fmt.Errorf("failed to marshal cached message: %w", err)
	}

	key := fmt.Sprintf("message:%s", messageID)
	expiration := 30 * 24 * time.Hour // 30 days

	return r.redis.Set(ctx, key, string(data), expiration)
}

// GetCachedMessage retrieves a cached message from Redis
func (r *messageCacheRepository) GetCachedMessage(ctx context.Context, messageID string) (*CachedMessage, error) {
	key := fmt.Sprintf("message:%s", messageID)

	data, err := r.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var cached CachedMessage
	if err := json.Unmarshal([]byte(data), &cached); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached message: %w", err)
	}

	return &cached, nil
}

// IsCached checks if a message is already cached
func (r *messageCacheRepository) IsCached(ctx context.Context, messageID string) (bool, error) {
	key := fmt.Sprintf("message:%s", messageID)

	count, err := r.redis.Exists(ctx, key)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
