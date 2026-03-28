package rate

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string, limit int, window time.Duration) (allowed bool, remaining int, reset time.Time, err error)
}

type RedisLimiter struct {
	client *redis.Client
	script *redis.Script
	prefix string
}

func NewRedisLimiter(client *redis.Client, prefix string) *RedisLimiter {
	if prefix == "" {
		prefix = "portlyn:ratelimit"
	}
	return &RedisLimiter{
		client: client,
		prefix: prefix,
		script: redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("PEXPIRE", KEYS[1], ARGV[1])
end
local ttl = redis.call("PTTL", KEYS[1])
return {current, ttl}
`),
	}
}

func (r *RedisLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, int, time.Time, error) {
	if limit <= 0 || window <= 0 {
		return true, 0, time.Now().UTC(), nil
	}
	values, err := r.script.Run(ctx, r.client, []string{r.key(key)}, window.Milliseconds()).Slice()
	if err != nil {
		return false, 0, time.Time{}, err
	}
	if len(values) != 2 {
		return false, 0, time.Time{}, errors.New("unexpected redis limiter response")
	}

	count := toInt64(values[0])
	ttlMillis := toInt64(values[1])
	if ttlMillis < 0 {
		ttlMillis = window.Milliseconds()
	}

	remaining := limit - int(count)
	if remaining < 0 {
		remaining = 0
	}
	reset := time.Now().UTC().Add(time.Duration(ttlMillis) * time.Millisecond)
	return count <= int64(limit), remaining, reset, nil
}

func (r *RedisLimiter) key(key string) string {
	return r.prefix + ":" + key
}

func toInt64(value any) int64 {
	switch typed := value.(type) {
	case int64:
		return typed
	case int:
		return int64(typed)
	case uint64:
		return int64(typed)
	case string:
		var out int64
		for _, ch := range typed {
			if ch < '0' || ch > '9' {
				break
			}
			out = out*10 + int64(ch-'0')
		}
		return out
	default:
		return 0
	}
}
