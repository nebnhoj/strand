package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/nebnhoj/strand/pkg/cache"
)

type redisCache struct {
	client *goredis.Client
}

func NewCache(client *goredis.Client) cache.Cache {
	return &redisCache{client: client}
}

func (c *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err == goredis.Nil {
		return nil, cache.ErrCacheMiss
	}
	return val, err
}

func (c *redisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

// DeleteByPattern deletes all keys matching the given glob pattern using SCAN + DEL.
func (c *redisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, next, err := c.client.Scan(ctx, cursor, pattern, 100).Result()
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

func (c *redisCache) Flush(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}
