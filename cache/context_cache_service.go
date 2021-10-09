package cache

import (
	"context"
	"time"
)

type ContextCacheService interface {
	Put(ctx context.Context, key string, obj interface{}, timeToLive time.Duration) error
	Expire(ctx context.Context, key string, timeToLive time.Duration) (bool, error)
	Get(ctx context.Context, key string) (interface{}, error)
	ContainsKey(ctx context.Context, key string) (bool, error)
	Remove(ctx context.Context, key string) (bool, error)
	Clear(ctx context.Context) error
	GetMany(ctx context.Context, keys []string) (map[string]interface{}, []string, error)
	GetManyStrings(ctx context.Context, keys []string) (map[string]string, []string, error)
	Keys(ctx context.Context) ([]string, error)
	Count(ctx context.Context) (int64, error)
	Size(ctx context.Context) (int64, error)
}
