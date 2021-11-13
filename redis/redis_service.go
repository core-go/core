package redis

import (
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

func (c *RedisService) Put(key string, obj interface{}, timeToLive time.Duration) error {
	return Set(c.Pool, key, obj, timeToLive)
}

func (c *RedisService) Expire(key string, timeToLive time.Duration) (bool, error) {
	return Expire(c.Pool, key, timeToLive)
}

func (c *RedisService) Get(key string) (interface{}, error) {
	return Get(c.Pool, key)
}

func (c *RedisService) ContainsKey(key string) (bool, error) {
	return Exists(c.Pool, key)
}

func (c *RedisService) Remove(key string) (bool, error) {
	return Delete(c.Pool, key)
}

func (c *RedisService) Clear() error {
	return Clear(c.Pool)
}

func (c *RedisService) GetMany(keys []string) (map[string]interface{}, []string, error) {
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

func (c *RedisService) Keys() ([]string, error) {
	return Keys(c.Pool)
}

func (c *RedisService) Count() (int64, error) {
	return Count(c.Pool)
}

func (c *RedisService) Size() (int64, error) {
	return Size(c.Pool)
}
