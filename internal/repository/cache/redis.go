package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func Connect(ctx context.Context, redisURL string) (*RedisCache, error) {
	if redisURL == "" {
		return nil, fmt.Errorf("empty redis URL")
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	log.Printf("[Cache] Connected to Redis 7.2 cluster at %s", opts.Addr)
	return &RedisCache{client: client}, nil
}

func (c *RedisCache) Client() *redis.Client {
	return c.client
}

func (c *RedisCache) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("redis client is nil")
	}
	return c.client.Get(ctx, key).Result()
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if c.client == nil {
		return fmt.Errorf("redis client is nil")
	}
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *RedisCache) Ping(ctx context.Context) error {
	if c.client == nil {
		return fmt.Errorf("redis client is nil")
	}
	return c.client.Ping(ctx).Err()
}

