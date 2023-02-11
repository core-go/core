package v9

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type SimpleRedisService struct {
	Client *redis.Client
}
func NewSimpleFailoverServiceByConfig(c FailoverConfig) *SimpleRedisService {
	client := NewFailoverClient(c)
	return &SimpleRedisService{client}
}
func NewSimpleRedisAdapterByConfig(c Config) (*SimpleRedisService, error) {
	return NewSimpleRedisServiceByConfig(c)
}
func NewSimpleRedisAdapter(redisUrl string) (*SimpleRedisService, error) {
	return NewSimpleRedisService(redisUrl)
}
func NewSimpleRedisServiceByConfig(c Config) (*SimpleRedisService, error) {
	client, err := NewRedisClientByConfig(c)
	if err != nil {
		return nil, err
	}
	return &SimpleRedisService{client}, nil
}
func NewSimpleRedisService(redisUrl string) (*SimpleRedisService, error) {
	client, err := NewRedisClient(redisUrl)
	if err != nil {
		return nil, err
	}
	return &SimpleRedisService{client}, nil
}

func (c *SimpleRedisService) Put(key string, obj interface{}, timeToLive time.Duration) error {
	return Set(context.TODO(), c.Client, key, obj, timeToLive)
}

func (c *SimpleRedisService) Expire(key string, timeToLive time.Duration) (bool, error) {
	return Expire(context.TODO(), c.Client, key, timeToLive)
}

func (c *SimpleRedisService) Get(key string) (interface{}, error) {
	return Get(context.TODO(), c.Client, key)
}

func (c *SimpleRedisService) ContainsKey(key string) (bool, error) {
	return Exists(context.TODO(), c.Client, key)
}

func (c *SimpleRedisService) Remove(key string) (bool, error) {
	return Delete(context.TODO(), c.Client, key)
}

func (c *SimpleRedisService) Clear() error {
	return Clear(context.TODO(), c.Client)
}

func (c *SimpleRedisService) GetMany(keys []string) (map[string]interface{}, []string, error) {
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

func (c *SimpleRedisService) GetManyStrings(keys []string) (map[string]string, []string, error) {
	return GetMany(context.TODO(), c.Client, keys)
}

func (c *SimpleRedisService) Keys() ([]string, error) {
	return Keys(context.TODO(), c.Client)
}

func (c *SimpleRedisService) Count() (int64, error) {
	return Count(context.TODO(), c.Client)
}

func (c *SimpleRedisService) Size() (int64, error) {
	return Size(context.TODO(), c.Client)
}
