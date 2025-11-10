package repository

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
	"github.com/srcndev/message-service/pkg/redis"
	"github.com/stretchr/testify/assert"
)

func setupMiniRedis(t *testing.T) (*miniredis.Miniredis, redis.Client) {
	mr := miniredis.RunT(t)

	rdb := goredis.NewClient(&goredis.Options{
		Addr: mr.Addr(),
	})

	// Wrap with custom client
	client := &testRedisClient{rdb: rdb}
	return mr, client
}

// testRedisClient implements redis.Client for testing
type testRedisClient struct {
	rdb *goredis.Client
}

func (c *testRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

func (c *testRedisClient) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *testRedisClient) Del(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

func (c *testRedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.rdb.Exists(ctx, keys...).Result()
}

func (c *testRedisClient) Close() error {
	return c.rdb.Close()
}

func (c *testRedisClient) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func TestMessageCacheRepository_CacheSentMessage_Success(t *testing.T) {
	mr, client := setupMiniRedis(t)
	defer mr.Close()

	repo := NewMessageCacheRepository(client)

	messageID := "test-message-id-123"
	sentAt := time.Now()

	err := repo.CacheSentMessage(context.Background(), messageID, sentAt)

	assert.NoError(t, err)

	// Verify in miniredis
	key := "message:" + messageID
	assert.True(t, mr.Exists(key))

	value, _ := mr.Get(key)
	var cached CachedMessage
	json.Unmarshal([]byte(value), &cached)
	assert.Equal(t, messageID, cached.MessageID)
}

func TestMessageCacheRepository_GetCachedMessage_Success(t *testing.T) {
	mr, client := setupMiniRedis(t)
	defer mr.Close()

	repo := NewMessageCacheRepository(client)

	messageID := "test-message-id-456"
	sentAt := time.Now()

	// First cache it
	_ = repo.CacheSentMessage(context.Background(), messageID, sentAt)

	// Then retrieve it
	cached, err := repo.GetCachedMessage(context.Background(), messageID)

	assert.NoError(t, err)
	assert.NotNil(t, cached)
	assert.Equal(t, messageID, cached.MessageID)
	assert.WithinDuration(t, sentAt, cached.SentAt, time.Second)
}

func TestMessageCacheRepository_GetCachedMessage_NotFound(t *testing.T) {
	_, client := setupMiniRedis(t)

	repo := NewMessageCacheRepository(client)

	messageID := "non-existent-id"

	cached, err := repo.GetCachedMessage(context.Background(), messageID)

	assert.Error(t, err)
	assert.Nil(t, cached)
}

func TestMessageCacheRepository_IsCached_True(t *testing.T) {
	mr, client := setupMiniRedis(t)
	defer mr.Close()

	repo := NewMessageCacheRepository(client)

	messageID := "cached-message-id"
	sentAt := time.Now()

	// Cache the message
	_ = repo.CacheSentMessage(context.Background(), messageID, sentAt)

	// Check if cached
	isCached, err := repo.IsCached(context.Background(), messageID)

	assert.NoError(t, err)
	assert.True(t, isCached)
}

func TestMessageCacheRepository_IsCached_False(t *testing.T) {
	_, client := setupMiniRedis(t)

	repo := NewMessageCacheRepository(client)

	messageID := "non-cached-message-id"

	isCached, err := repo.IsCached(context.Background(), messageID)

	assert.NoError(t, err)
	assert.False(t, isCached)
}

func TestMessageCacheRepository_CacheSentMessage_TTL(t *testing.T) {
	mr, client := setupMiniRedis(t)
	defer mr.Close()

	repo := NewMessageCacheRepository(client)

	messageID := "ttl-test-message-id"
	sentAt := time.Now()

	err := repo.CacheSentMessage(context.Background(), messageID, sentAt)
	assert.NoError(t, err)

	// Check TTL in miniredis
	key := "message:" + messageID
	ttl := mr.TTL(key)
	assert.True(t, ttl > 0, "TTL should be set")

	// Should be approximately 30 days (allow 1 second difference)
	expectedTTL := 30 * 24 * time.Hour
	assert.InDelta(t, expectedTTL.Seconds(), ttl.Seconds(), 1.0)
}

func TestMessageCacheRepository_CacheSentMessage_KeyFormat(t *testing.T) {
	mr, client := setupMiniRedis(t)
	defer mr.Close()

	repo := NewMessageCacheRepository(client)

	messageID := "key-format-test-id"
	sentAt := time.Now()

	err := repo.CacheSentMessage(context.Background(), messageID, sentAt)
	assert.NoError(t, err)

	// Verify key format: message:{messageId}
	expectedKey := "message:" + messageID
	assert.True(t, mr.Exists(expectedKey))
}

func TestMessageCacheRepository_GetCachedMessage_InvalidJSON(t *testing.T) {
	mr, client := setupMiniRedis(t)
	defer mr.Close()

	repo := NewMessageCacheRepository(client)

	messageID := "invalid-json-id"
	key := "message:" + messageID

	// Set invalid JSON directly in miniredis
	mr.Set(key, "invalid json data")

	cached, err := repo.GetCachedMessage(context.Background(), messageID)

	assert.Error(t, err)
	assert.Nil(t, cached)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

func TestMessageCacheRepository_MultipleCachedMessages(t *testing.T) {
	mr, client := setupMiniRedis(t)
	defer mr.Close()

	repo := NewMessageCacheRepository(client)

	messages := []struct {
		id     string
		sentAt time.Time
	}{
		{"msg-1", time.Now().Add(-1 * time.Hour)},
		{"msg-2", time.Now().Add(-2 * time.Hour)},
		{"msg-3", time.Now()},
	}

	// Cache all messages
	for _, msg := range messages {
		err := repo.CacheSentMessage(context.Background(), msg.id, msg.sentAt)
		assert.NoError(t, err)
	}

	// Verify all are cached
	for _, msg := range messages {
		isCached, err := repo.IsCached(context.Background(), msg.id)
		assert.NoError(t, err)
		assert.True(t, isCached)

		cached, err := repo.GetCachedMessage(context.Background(), msg.id)
		assert.NoError(t, err)
		assert.Equal(t, msg.id, cached.MessageID)
		assert.WithinDuration(t, msg.sentAt, cached.SentAt, time.Second)
	}
}

func TestMessageCacheRepository_InterfaceCompliance(t *testing.T) {
	var _ MessageCacheRepository = (*messageCacheRepository)(nil)

	_, client := setupMiniRedis(t)
	repo := NewMessageCacheRepository(client)
	assert.NotNil(t, repo)
}
