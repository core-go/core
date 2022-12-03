package v8

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisService struct {
	Client *redis.Client
}
func NewRedisServiceByConfig(c Config) (*RedisService, error) {
	client, err := NewRedisClientByConfig(c)
	if err != nil {
		return nil, err
	}
	return &RedisService{client}, nil
}
func NewRedisService(redisUrl string) (*RedisService, error) {
	client, err := NewRedisClient(redisUrl)
	if err != nil {
		return nil, err
	}
	return &RedisService{client}, nil
}

func (c *RedisService) Put(ctx context.Context, key string, obj interface{}, timeToLive time.Duration) error {
	return Set(ctx, c.Client, key, obj, timeToLive)
}

func (c *RedisService) Expire(ctx context.Context, key string, timeToLive time.Duration) (bool, error) {
	return Expire(ctx, c.Client, key, timeToLive)
}

func (c *RedisService) Get(ctx context.Context, key string) (interface{}, error) {
	return Get(ctx, c.Client, key)
}

func (c *RedisService) ContainsKey(ctx context.Context, key string) (bool, error) {
	return Exists(ctx, c.Client, key)
}

func (c *RedisService) Remove(ctx context.Context, key string) (bool, error) {
	return Delete(ctx, c.Client, key)
}

func (c *RedisService) Clear(ctx context.Context, ) error {
	return Clear(ctx, c.Client)
}

func (c *RedisService) GetMany(ctx context.Context, keys []string) (map[string]interface{}, []string, error) {
	m2 := make(map[string]interface{})
	m, n, err := GetMany(ctx, c.Client, keys)
	if err != nil {
		return m2, n, err
	}
	for k, v := range m {
		m2[k] = v
	}
	return m2, n, nil
}

func (c *RedisService) GetManyStrings(ctx context.Context, keys []string) (map[string]string, []string, error) {
	return GetMany(ctx, c.Client, keys)
}

func (c *RedisService) Keys(ctx context.Context, ) ([]string, error) {
	return Keys(ctx, c.Client)
}

func (c *RedisService) Count(ctx context.Context) (int64, error) {
	return Count(ctx, c.Client)
}

func (c *RedisService) Size(ctx context.Context) (int64, error) {
	return Size(ctx, c.Client)
}
