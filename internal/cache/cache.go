package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	rc  *redis.Client
	ttl time.Duration
}

func New(rc *redis.Client, ttl time.Duration) *Cache {
	return &Cache{
		rc:  rc,
		ttl: ttl,
	}
}

func (c *Cache) Get(ctx context.Context, key string) (string, bool, error) {
	val, err := c.rc.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (c *Cache) Set(ctx context.Context, key, value string) error {
	return c.rc.Set(ctx, key, value, c.ttl).Err()
}
