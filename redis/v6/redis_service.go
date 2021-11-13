package v6

import (
	"github.com/go-redis/redis"
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

func (c *RedisService) Put(key string, obj interface{}, timeToLive time.Duration) error {
	return Set(c.Client, key, obj, timeToLive)
}

func (c *RedisService) Expire(key string, timeToLive time.Duration) (bool, error) {
	return Expire(c.Client, key, timeToLive)
}

func (c *RedisService) Get(key string) (interface{}, error) {
	return Get(c.Client, key)
}

func (c *RedisService) ContainsKey(key string) (bool, error) {
	return Exists(c.Client, key)
}

func (c *RedisService) Remove(key string) (bool, error) {
	return Delete(c.Client, key)
}

func (c *RedisService) Clear() error {
	return Clear(c.Client)
}

func (c *RedisService) GetMany(keys []string) (map[string]interface{}, []string, error) {
	m2 := make(map[string]interface{})
	m, n, err := GetMany(c.Client, keys)
	if err != nil {
		return m2, n, err
	}
	for k, v := range m {
		m2[k] = v
	}
	return m2, n, nil
}

func (c *RedisService) GetManyStrings(keys []string) (map[string]string, []string, error) {
	return GetMany(c.Client, keys)
}

func (c *RedisService) Keys() ([]string, error) {
	return Keys(c.Client)
}

func (c *RedisService) Count() (int64, error) {
	return Count(c.Client)
}

func (c *RedisService) Size() (int64, error) {
	return Size(c.Client)
}
