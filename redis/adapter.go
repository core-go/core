package redis

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisAdapter struct {
	Pool *redis.Pool
}

func NewRedisAdapterByConfig(c Config) (*RedisAdapter, error) {
	pool, err := NewRedisPoolByConfig(c)
	if err != nil {
		return nil, err
	}
	return &RedisAdapter{pool}, nil
}
func NewRedisAdapter(redisUrl string) (*RedisAdapter, error) {
	pool, err := NewRedisPool(redisUrl)
	if err != nil {
		return nil, err
	}
	return &RedisAdapter{pool}, nil
}

func (c *RedisAdapter) Put(ctx context.Context, key string, obj interface{}, timeToLive time.Duration) error {
	return Set(c.Pool, key, obj, timeToLive)
}

func (c *RedisAdapter) Expire(ctx context.Context, key string, timeToLive time.Duration) (bool, error) {
	return Expire(c.Pool, key, timeToLive)
}

func (c *RedisAdapter) Get(ctx context.Context, key string) (string, error) {
	return Get(c.Pool, key)
}

func (c *RedisAdapter) ContainsKey(ctx context.Context, key string) (bool, error) {
	return Exists(c.Pool, key)
}

func (c *RedisAdapter) Remove(ctx context.Context, key string) (bool, error) {
	return Delete(c.Pool, key)
}

func (c *RedisAdapter) Clear(ctx context.Context) error {
	return Clear(c.Pool)
}

func (c *RedisAdapter) GetMany(keys []string) (map[string]string, []string, error) {
	return GetMany(c.Pool, keys)
}

func (c *RedisAdapter) Keys(ctx context.Context) ([]string, error) {
	return Keys(c.Pool)
}

func (c *RedisAdapter) Count(ctx context.Context) (int64, error) {
	return Count(c.Pool)
}

func (c *RedisAdapter) Size(ctx context.Context) (int64, error) {
	return Size(c.Pool)
}
