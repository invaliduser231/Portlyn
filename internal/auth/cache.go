package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"portlyn/internal/domain"
)

type CachedAuthEntry struct {
	User      domain.User `json:"user"`
	GroupIDs  []uint      `json:"group_ids"`
	ExpiresAt time.Time   `json:"expires_at"`
}

type AuthCache interface {
	Get(ctx context.Context, tokenHash string) (*CachedAuthEntry, bool, error)
	Set(ctx context.Context, tokenHash string, entry CachedAuthEntry, ttl time.Duration) error
	InvalidateToken(ctx context.Context, tokenHash string) error
	InvalidateUser(ctx context.Context, userID uint) error
	Flush(ctx context.Context) error
}

type RedisAuthCache struct {
	client      *redis.Client
	entryPrefix string
	userPrefix  string
}

func NewRedisAuthCache(client *redis.Client, prefix string) *RedisAuthCache {
	if prefix == "" {
		prefix = "portlyn:authcache"
	}
	return &RedisAuthCache{
		client:      client,
		entryPrefix: prefix + ":entry",
		userPrefix:  prefix + ":user",
	}
}

func (c *RedisAuthCache) Get(ctx context.Context, tokenHash string) (*CachedAuthEntry, bool, error) {
	value, err := c.client.Get(ctx, c.entryKey(tokenHash)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, false, err
	}

	var entry CachedAuthEntry
	if err := json.Unmarshal(value, &entry); err != nil {
		return nil, false, err
	}
	if time.Now().UTC().After(entry.ExpiresAt) {
		_ = c.InvalidateToken(ctx, tokenHash)
		return nil, false, nil
	}
	return &entry, true, nil
}

func (c *RedisAuthCache) Set(ctx context.Context, tokenHash string, entry CachedAuthEntry, ttl time.Duration) error {
	payload, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	pipe := c.client.TxPipeline()
	pipe.Set(ctx, c.entryKey(tokenHash), payload, ttl)
	pipe.SAdd(ctx, c.userKey(entry.User.ID), tokenHash)
	pipe.Expire(ctx, c.userKey(entry.User.ID), ttl)
	_, err = pipe.Exec(ctx)
	return err
}

func (c *RedisAuthCache) InvalidateToken(ctx context.Context, tokenHash string) error {
	return c.client.Del(ctx, c.entryKey(tokenHash)).Err()
}

func (c *RedisAuthCache) InvalidateUser(ctx context.Context, userID uint) error {
	key := c.userKey(userID)
	tokenHashes, err := c.client.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	pipe := c.client.TxPipeline()
	for _, tokenHash := range tokenHashes {
		pipe.Del(ctx, c.entryKey(tokenHash))
	}
	pipe.Del(ctx, key)
	_, err = pipe.Exec(ctx)
	return err
}

func (c *RedisAuthCache) Flush(ctx context.Context) error {
	var cursor uint64
	for {
		keys, next, err := c.client.Scan(ctx, cursor, c.entryPrefix+":*", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return nil
}

func (c *RedisAuthCache) entryKey(tokenHash string) string {
	return fmt.Sprintf("%s:%s", c.entryPrefix, tokenHash)
}

func (c *RedisAuthCache) userKey(userID uint) string {
	return fmt.Sprintf("%s:%d", c.userPrefix, userID)
}
