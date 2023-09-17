package echo

import "time"

type CacheService interface {
	Put(key string, obj interface{}, timeToLive time.Duration) error
	GetManyStrings(key []string) (map[string]string, []string, error)
}
