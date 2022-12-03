package v8

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisAdapter struct {
	Client *redis.Client
}
func NewRedisAdapterByConfig(c Config) (*RedisAdapter, error) {
	client, err := NewRedisClientByConfig(c)
	if err != nil {
		return nil, err
	}
	return &RedisAdapter{client}, nil
}
func NewRedisAdapter(redisUrl string) (*RedisAdapter, error) {
	client, err := NewRedisClient(redisUrl)
	if err != nil {
		return nil, err
	}
	return &RedisAdapter{client}, nil
}

func (c *RedisAdapter) Put(ctx context.Context, key string, obj interface{}, timeToLive time.Duration) error {
	return Set(ctx, c.Client, key, obj, timeToLive)
}

func (c *RedisAdapter) Expire(ctx context.Context, key string, timeToLive time.Duration) (bool, error) {
	return Expire(ctx, c.Client, key, timeToLive)
}

func (c *RedisAdapter) Get(ctx context.Context, key string) (string, error) {
	return Get(ctx, c.Client, key)
}

func (c *RedisAdapter) ContainsKey(ctx context.Context, key string) (bool, error) {
	return Exists(ctx, c.Client, key)
}

func (c *RedisAdapter) Remove(ctx context.Context, key string) (bool, error) {
	return Delete(ctx, c.Client, key)
}

func (c *RedisAdapter) Clear(ctx context.Context, ) error {
	return Clear(ctx, c.Client)
}

func (c *RedisAdapter) GetMany(ctx context.Context, keys []string) (map[string]string, []string, error) {
	return GetMany(ctx, c.Client, keys)
}

func (c *RedisAdapter) Keys(ctx context.Context, ) ([]string, error) {
	return Keys(ctx, c.Client)
}

func (c *RedisAdapter) Count(ctx context.Context) (int64, error) {
	return Count(ctx, c.Client)
}

func (c *RedisAdapter) Size(ctx context.Context) (int64, error) {
	return Size(ctx, c.Client)
}
