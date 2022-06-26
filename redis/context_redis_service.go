package redis

import (
	"context"
	"github.com/garyburd/redigo/redis"
	"time"
)

type ContextRedisService struct {
	Pool *redis.Pool
}
func NewContextRedisAdapterByConfig(c Config) (*ContextRedisService, error) {
	return NewContextRedisServiceByConfig(c)
}
func NewContextRedisAdapter(redisUrl string) (*ContextRedisService, error) {
	return NewContextRedisService(redisUrl)
}
func NewContextRedisServiceByConfig(c Config) (*ContextRedisService, error) {
	pool, err := NewRedisPoolByConfig(c)
	if err != nil {
		return nil, err
	}
	return &ContextRedisService{pool}, nil
}
func NewContextRedisService(redisUrl string) (*ContextRedisService, error) {
	pool, err := NewRedisPool(redisUrl)
	if err != nil {
		return nil, err
	}
	return &ContextRedisService{pool}, nil
}

func (c *ContextRedisService) Put(ctx context.Context, key string, obj interface{}, timeToLive time.Duration) error {
	return Set(c.Pool, key, obj, timeToLive)
}

func (c *ContextRedisService) Expire(ctx context.Context, key string, timeToLive time.Duration) (bool, error) {
	return Expire(c.Pool, key, timeToLive)
}

func (c *ContextRedisService) Get(ctx context.Context, key string) (interface{}, error) {
	return Get(c.Pool, key)
}

func (c *ContextRedisService) ContainsKey(ctx context.Context, key string) (bool, error) {
	return Exists(c.Pool, key)
}

func (c *ContextRedisService) Remove(ctx context.Context, key string) (bool, error) {
	return Delete(c.Pool, key)
}

func (c *ContextRedisService) Clear(ctx context.Context) error {
	return Clear(c.Pool)
}

func (c *ContextRedisService) GetMany(ctx context.Context, keys []string) (map[string]interface{}, []string, error) {
	m2 := make(map[string]interface{})
	m, n, err := GetMany(c.Pool, keys)
	if err != nil {
		return m2, n, err
	}
	for k, v := range m {
		m2[k] = v
	}
	return m2, n, nil
}

func (c *ContextRedisService) GetManyStrings(keys []string) (map[string]string, []string, error) {
	return GetMany(c.Pool, keys)
}

func (c *ContextRedisService) Keys(ctx context.Context) ([]string, error) {
	return Keys(c.Pool)
}

func (c *ContextRedisService) Count(ctx context.Context) (int64, error) {
	return Count(c.Pool)
}

func (c *ContextRedisService) Size(ctx context.Context) (int64, error) {
	return Size(c.Pool)
}
