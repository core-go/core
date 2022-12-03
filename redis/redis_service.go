package redis

import (
	"context"
	"github.com/garyburd/redigo/redis"
	"time"
)

type RedisService struct {
	Pool *redis.Pool
}
func NewRedisServiceByConfig(c Config) (*RedisService, error) {
	pool, err := NewRedisPoolByConfig(c)
	if err != nil {
		return nil, err
	}
	return &RedisService{pool}, nil
}
func NewRedisService(redisUrl string) (*RedisService, error) {
	pool, err := NewRedisPool(redisUrl)
	if err != nil {
		return nil, err
	}
	return &RedisService{pool}, nil
}

func (c *RedisService) Put(ctx context.Context, key string, obj interface{}, timeToLive time.Duration) error {
	return Set(c.Pool, key, obj, timeToLive)
}

func (c *RedisService) Expire(ctx context.Context, key string, timeToLive time.Duration) (bool, error) {
	return Expire(c.Pool, key, timeToLive)
}

func (c *RedisService) Get(ctx context.Context, key string) (interface{}, error) {
	return Get(c.Pool, key)
}

func (c *RedisService) ContainsKey(ctx context.Context, key string) (bool, error) {
	return Exists(c.Pool, key)
}

func (c *RedisService) Remove(ctx context.Context, key string) (bool, error) {
	return Delete(c.Pool, key)
}

func (c *RedisService) Clear(ctx context.Context) error {
	return Clear(c.Pool)
}

func (c *RedisService) GetMany(ctx context.Context, keys []string) (map[string]interface{}, []string, error) {
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

func (c *RedisService) GetManyStrings(keys []string) (map[string]string, []string, error) {
	return GetMany(c.Pool, keys)
}

func (c *RedisService) Keys(ctx context.Context) ([]string, error) {
	return Keys(c.Pool)
}

func (c *RedisService) Count(ctx context.Context) (int64, error) {
	return Count(c.Pool)
}

func (c *RedisService) Size(ctx context.Context) (int64, error) {
	return Size(c.Pool)
}
