package cache

import "time"

type CacheService interface {
	Put(key string, obj interface{}, timeToLive time.Duration) error
	Expire(key string, timeToLive time.Duration) (bool, error)
	Get(key string) (interface{}, error)
	ContainsKey(key string) (bool, error)
	Remove(key string) (bool, error)
	Clear() error
	GetMany(keys []string) (map[string]interface{}, []string, error)
	GetManyStrings(keys []string) (map[string]string, []string, error)
	Keys() ([]string, error)
	Count() (int64, error)
	Size() (int64, error)
}
