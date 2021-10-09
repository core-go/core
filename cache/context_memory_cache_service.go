package cache

import (
	"context"
	"errors"
	"time"
	"unsafe"
)

// ContextMemoryCacheService manage all custom caching action
type ContextMemoryCacheService struct {
	client *Client
	close  chan struct{}
}

func NewContextMemoryCacheServiceByConfig(conf CacheConfig) (*ContextMemoryCacheService, error) {
	return NewContextMemoryCacheService(conf.Size, conf.CleaningEnable, conf.CleaningInterval)
}
// NewContextMemoryCacheService init new instance
func NewContextMemoryCacheService(size int64, cleaningEnable bool, cleaningInterval time.Duration) (*ContextMemoryCacheService, error) {
	currentSession := &ContextMemoryCacheService{NewClient(size, cleaningEnable), make(chan struct{})}

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
func (c *ContextMemoryCacheService) Get(ctx context.Context, key string) (interface{}, error) {
	obj, err := c.client.Read(key)
	if err != nil {
		return nil, err
	}

	item, ok := obj.(Item)
	if !ok {
		return nil, errors.New("can not map object to Item model")
	}

	if item.Expires < time.Now().UnixNano() {
		return nil, nil
	}

	return item.Data, nil
}

// Get return value based on the list of keys provided
func (c *ContextMemoryCacheService) GetMany(ctx context.Context, keys []string) (map[string]interface{}, []string, error) {
	var itemFound map[string]interface{}
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

// Get return value based on the list of keys provided
func (c *ContextMemoryCacheService) GetManyStrings(ctx context.Context, keys []string) (map[string]string, []string, error) {
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

		itemFound[key] = item.Data.(string)
	}

	return itemFound, itemNotFound, nil
}

// Get return value based on the key provided
func (c *ContextMemoryCacheService) ContainsKey(ctx context.Context, key string) (bool, error) {
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

// Put new record set key and value
func (c *ContextMemoryCacheService) Put(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	if expire == 0 {
		expire = 24 * time.Hour
	}

	if err := c.client.Push(key, Item{
		Data:    value,
		Expires: time.Now().Add(expire).UnixNano(),
	}); err != nil {
		return err
	}

	return nil
}

// Expire new value over the key provided
func (c *ContextMemoryCacheService) Expire(ctx context.Context, key string, expire time.Duration) (bool, error) {
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
func (c *ContextMemoryCacheService) Remove(ctx context.Context, key string) (bool, error) {
	if _, err := c.client.Get(key); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (c *ContextMemoryCacheService) Clear(ctx context.Context) error {
	return nil
}

// Count return number of records
func (c *ContextMemoryCacheService) Count(ctx context.Context) (int64, error) {
	return int64(c.client.GetNumberOfKeys()), nil
}

func (c *ContextMemoryCacheService) Keys(ctx context.Context) ([]string, error) {
	return c.client.Getkeys(), nil
}

// GetDBSize method return redis database size
func (c *ContextMemoryCacheService) Size(ctx context.Context) (int64, error) {
	return int64(unsafe.Sizeof(c.client)), nil
}

// Close closes the cache and frees up resources.
func (c *ContextMemoryCacheService) Close(ctx context.Context) error {
	c.close <- struct{}{}
	c.client = NewClient(10*1024*1024, true) // 10 * 1024 * 1024 for 10 mb
	return nil
}
