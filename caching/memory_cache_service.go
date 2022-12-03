package caching

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"unsafe"
)

type MemoryCacheService struct {
	client *Client
	close  chan struct{}
}
func NewCacheService(size int64, cleaningEnable bool, cleaningInterval time.Duration) (*MemoryCacheService, error) {
	return NewMemoryCacheService(size, cleaningEnable, cleaningInterval)
}
func NewCacheAdapterByConfig(conf CacheConfig) (*MemoryCacheService, error) {
	return NewMemoryCacheService(conf.Size, conf.CleaningEnable, conf.CleaningInterval)
}
func NewCacheAdapter(size int64, cleaningEnable bool, cleaningInterval time.Duration) (*MemoryCacheService, error) {
	return NewMemoryCacheService(size, cleaningEnable, cleaningInterval)
}
func NewMemoryCacheAdapterByConfig(conf CacheConfig) (*MemoryCacheService, error) {
	return NewMemoryCacheService(conf.Size, conf.CleaningEnable, conf.CleaningInterval)
}
func NewMemoryCacheAdapter(size int64, cleaningEnable bool, cleaningInterval time.Duration) (*MemoryCacheService, error) {
	return NewMemoryCacheService(size, cleaningEnable, cleaningInterval)
}
func NewMemoryCacheServiceByConfig(conf CacheConfig) (*MemoryCacheService, error) {
	return NewMemoryCacheService(conf.Size, conf.CleaningEnable, conf.CleaningInterval)
}
func NewMemoryCacheService(size int64, cleaningEnable bool, cleaningInterval time.Duration) (*MemoryCacheService, error) {
	currentSession := &MemoryCacheService{NewClient(size, cleaningEnable), make(chan struct{})}

	// Check record expiration time and remove
	go func() {
		ticker := time.NewTicker(cleaningInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				items := currentSession.client.GetItems()
				items.Range(func(key, value interface{}) bool {
					item := value.(Item)

					if item.Expires < time.Now().UnixNano() {
						k, _ := key.(string)
						currentSession.client.Get(k)
					}

					return true
				})

			case <-currentSession.close:
				return
			}
		}
	}()

	return currentSession, nil
}

// Get return value based on the key provided
func (c *MemoryCacheService) Get(ctx context.Context, key string) (string, error) {
	obj, err := c.client.Read(key)
	if err != nil {
		return "", err
	}

	item, ok := obj.(Item)
	if !ok {
		return "", errors.New("can not map object to Item model")
	}

	if item.Expires < time.Now().UnixNano() {
		return "", nil
	}

	return item.Data, nil
}
func (c *MemoryCacheService) GetMany(ctx context.Context, keys []string) (map[string]string, []string, error) {
	var itemFound map[string]string
	var itemNotFound []string

	for _, key := range keys {
		obj, err := c.client.Read(key)
		if obj == nil && err == nil {
			itemNotFound = append(itemNotFound, key)
		}

		item, ok := obj.(Item)
		if !ok {
			return nil, nil, errors.New("can not map object to Item model")
		}

		itemFound[key] = item.Data
	}

	return itemFound, itemNotFound, nil
}

func (c *MemoryCacheService) ContainsKey(ctx context.Context, key string) (bool, error) {
	obj, err := c.client.Read(key)
	if err != nil {
		return false, err
	}

	item, ok := obj.(Item)
	if !ok {
		return false, errors.New("can not map object to Item model")
	}

	if item.Expires < time.Now().UnixNano() {
		return false, nil
	}

	return true, nil
}

func (c *MemoryCacheService) Put(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	if expire == 0 {
		expire = 24 * time.Hour
	}
	var v string
	v, ok := value.(string)
	if ok == false {
		json, err := json.Marshal(value)
		if err != nil {
			return err
		}
		v = string(json)
	}
	if err := c.client.Push(key, Item{
		Data:    v,
		Expires: time.Now().Add(expire).UnixNano(),
	}); err != nil {
		return err
	}

	return nil
}

// Expire new value over the key provided
func (c *MemoryCacheService) Expire(ctx context.Context, key string, expire time.Duration) (bool, error) {
	val, err := c.client.Get(key)
	if err != nil {
		return false, err
	}

	if err := c.client.Push(key, Item{
		Data:    val,
		Expires: time.Now().Add(expire).UnixNano(),
	}); err != nil {
		return false, err
	}

	return true, nil
}

// Remove deletes the key and its value from the cache.
func (c *MemoryCacheService) Remove(ctx context.Context, key string) (bool, error) {
	if _, err := c.client.Get(key); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (c *MemoryCacheService) Clear(ctx context.Context) error {
	return nil
}

func (c *MemoryCacheService) Count(ctx context.Context) (int64, error) {
	return int64(c.client.GetNumberOfKeys()), nil
}

func (c *MemoryCacheService) Keys(ctx context.Context) ([]string, error) {
	return c.client.Getkeys(), nil
}

func (c *MemoryCacheService) Size(ctx context.Context) (int64, error) {
	return int64(unsafe.Sizeof(c.client)), nil
}

func (c *MemoryCacheService) Close(ctx context.Context) error {
	c.close <- struct{}{}
	c.client = NewClient(10*1024*1024, true) // 10 * 1024 * 1024 for 10 mb
	return nil
}
