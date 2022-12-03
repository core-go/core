package redis

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type SimpleRedisService struct {
	Pool *redis.Pool
}

func NewSimpleRedisAdapterByConfig(c Config) (*SimpleRedisService, error) {
	return NewSimpleRedisServiceByConfig(c)
}
func NewSimpleRedisAdapter(redisUrl string) (*SimpleRedisService, error) {
	return NewSimpleRedisService(redisUrl)
}
func NewSimpleRedisServiceByConfig(c Config) (*SimpleRedisService, error) {
	pool, err := NewRedisPoolByConfig(c)
	if err != nil {
		return nil, err
	}
	return &SimpleRedisService{pool}, nil
}
func NewSimpleRedisService(redisUrl string) (*SimpleRedisService, error) {
	pool, err := NewRedisPool(redisUrl)
	if err != nil {
		return nil, err
	}
	return &SimpleRedisService{pool}, nil
}

func (c *SimpleRedisService) Put(key string, obj interface{}, timeToLive time.Duration) error {
	return Set(c.Pool, key, obj, timeToLive)
}

func (c *SimpleRedisService) Expire(key string, timeToLive time.Duration) (bool, error) {
	return Expire(c.Pool, key, timeToLive)
}

func (c *SimpleRedisService) Get(key string) (interface{}, error) {
	return Get(c.Pool, key)
}

func (c *SimpleRedisService) ContainsKey(key string) (bool, error) {
	return Exists(c.Pool, key)
}

func (c *SimpleRedisService) Remove(key string) (bool, error) {
	return Delete(c.Pool, key)
}

func (c *SimpleRedisService) Clear() error {
	return Clear(c.Pool)
}

func (c *SimpleRedisService) GetMany(keys []string) (map[string]interface{}, []string, error) {
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

func (c *SimpleRedisService) GetManyStrings(keys []string) (map[string]string, []string, error) {
	return GetMany(c.Pool, keys)
}

func (c *SimpleRedisService) Keys() ([]string, error) {
	return Keys(c.Pool)
}

func (c *SimpleRedisService) Count() (int64, error) {
	return Count(c.Pool)
}

func (c *SimpleRedisService) Size() (int64, error) {
	return Size(c.Pool)
}
