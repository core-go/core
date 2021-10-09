package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type ContextRedisService struct {
	Client *redis.Client
}

func NewContextRedisServiceByConfig(c Config) (*ContextRedisService, error) {
	client, err := NewRedisClientByConfig(c)
	if err != nil {
		return nil, err
	}
	return &ContextRedisService{client}, nil
}
func NewContextRedisService(redisUrl string) (*ContextRedisService, error) {
	client, err := NewRedisClient(redisUrl)
	if err != nil {
		return nil, err
	}
	return &ContextRedisService{client}, nil
}

func (c *ContextRedisService) Put(ctx context.Context, key string, obj interface{}, timeToLive time.Duration) error {
	return Set(ctx, c.Client, key, obj, timeToLive)
}

func (c *ContextRedisService) Expire(ctx context.Context, key string, timeToLive time.Duration) (bool, error) {
	return Expire(ctx, c.Client, key, timeToLive)
}

func (c *ContextRedisService) Get(ctx context.Context, key string) (interface{}, error) {
	return Get(ctx, c.Client, key)
}

func (c *ContextRedisService) ContainsKey(ctx context.Context, key string) (bool, error) {
	return Exists(ctx, c.Client, key)
}

func (c *ContextRedisService) Remove(ctx context.Context, key string) (bool, error) {
	return Delete(ctx, c.Client, key)
}

func (c *ContextRedisService) Clear(ctx context.Context, ) error {
	return Clear(ctx, c.Client)
}

func (c *ContextRedisService) GetMany(ctx context.Context, keys []string) (map[string]interface{}, []string, error) {
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

func (c *ContextRedisService) GetManyStrings(ctx context.Context, keys []string) (map[string]string, []string, error) {
	return GetMany(ctx, c.Client, keys)
}

func (c *ContextRedisService) Keys(ctx context.Context, ) ([]string, error) {
	return Keys(ctx, c.Client)
}

func (c *ContextRedisService) Count(ctx context.Context) (int64, error) {
	return Count(ctx, c.Client)
}

func (c *ContextRedisService) Size(ctx context.Context) (int64, error) {
	return Size(ctx, c.Client)
}
