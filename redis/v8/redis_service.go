package v8

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisService struct {
	Client *redis.Client
}
func NewRedisAdapterByConfig(c Config) (*RedisService, error) {
	return NewRedisServiceByConfig(c)
}
func NewRedisAdapter(redisUrl string) (*RedisService, error) {
	return NewRedisService(redisUrl)
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
	return Set(context.TODO(), c.Client, key, obj, timeToLive)
}

func (c *RedisService) Expire(key string, timeToLive time.Duration) (bool, error) {
	return Expire(context.TODO(), c.Client, key, timeToLive)
}

func (c *RedisService) Get(key string) (interface{}, error) {
	return Get(context.TODO(), c.Client, key)
}

func (c *RedisService) ContainsKey(key string) (bool, error) {
	return Exists(context.TODO(), c.Client, key)
}

func (c *RedisService) Remove(key string) (bool, error) {
	return Delete(context.TODO(), c.Client, key)
}

func (c *RedisService) Clear() error {
	return Clear(context.TODO(), c.Client)
}

func (c *RedisService) GetMany(keys []string) (map[string]interface{}, []string, error) {
	m2 := make(map[string]interface{})
	m, n, err := GetMany(context.TODO(), c.Client, keys)
	if err != nil {
		return m2, n, err
	}
	for k, v := range m {
		m2[k] = v
	}
	return m2, n, nil
}

func (c *RedisService) GetManyStrings(keys []string) (map[string]string, []string, error) {
	return GetMany(context.TODO(), c.Client, keys)
}

func (c *RedisService) Keys() ([]string, error) {
	return Keys(context.TODO(), c.Client)
}

func (c *RedisService) Count() (int64, error) {
	return Count(context.TODO(), c.Client)
}

func (c *RedisService) Size() (int64, error) {
	return Size(context.TODO(), c.Client)
}
