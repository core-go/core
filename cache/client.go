package cache

import (
	"errors"
	"sync"
	"unsafe"
)

// LinearClient contains all attributes
type Client struct {
	items              sync.Map
	keys               []string
	cleaningEnable     bool
	linearSizes        int64
	linearCurrentSizes int64
	rwMutex            sync.RWMutex
}

// NewClient return new instance
func NewClient(linearSizes int64, cleaningEnable bool) *Client {
	currentLinear := Client{
		keys:               []string{},
		items:              sync.Map{},
		cleaningEnable:     cleaningEnable,
		linearSizes:        linearSizes,
		linearCurrentSizes: 0,
		rwMutex:            sync.RWMutex{},
	}

	return &currentLinear
}

// Push item to the linear by key
func (c *Client) Push(key string, value interface{}) error {
	itemSize := int64(unsafe.Sizeof(key)) + int64(unsafe.Sizeof(value))
	if itemSize > c.linearSizes || key == "" {
		return errors.New("key is empty or linear not enough space")
	}

	// Clean space for new item
	if c.cleaningEnable {
		for c.linearCurrentSizes+itemSize > c.linearSizes {
			c.Take()
		}
	}

	c.rwMutex.Lock()
	c.items.LoadOrStore(key, value)
	c.linearCurrentSizes += int64(unsafe.Sizeof(value))
	c.keys = append(c.keys, key)
	c.rwMutex.Unlock()

	return nil
}

// Pop return the last item from the linear and remove it
func (c *Client) Pop() (interface{}, error) {
	if c.IsEmpty() {
		return nil, errors.New("the linear is empty")
	}

	lastItemIndex := len(c.keys) - 1
	item, exits := c.items.Load(c.keys[lastItemIndex])
	if !exits {
		return nil, nil
	}

	c.rwMutex.Lock()
	c.items.Delete(c.keys[lastItemIndex])
	c.linearCurrentSizes -= int64(unsafe.Sizeof(item))
	c.keys = removeItemByIndex(c.keys, lastItemIndex) //Update keys slice after remove that key from items map
	c.rwMutex.Unlock()

	return item, nil
}

// Take return the first item from the linear and remove it
func (c *Client) Take() (interface{}, error) {
	if c.IsEmpty() {
		return nil, errors.New("the linear is empty")
	}

	c.rwMutex.Lock()
	item, exits := c.items.Load(c.keys[0])
	if !exits {
		c.rwMutex.Unlock()
		return nil, nil
	}

	c.items.Delete(c.keys[0])
	c.linearCurrentSizes -= int64(unsafe.Sizeof(item))
	c.keys = removeItemByIndex(c.keys, 0) //Update keys slice after remove that key from items map
	c.rwMutex.Unlock()

	return item, nil
}

// Get method return the item by key from linear and remove it
// Goroutine: https://stackoverflow.com/questions/20945069/catching-return-values-from-goroutines
func (c *Client) Get(key string) (interface{}, error) {
	if c.IsEmpty() {
		return nil, errors.New("linear is empty")
	}

	var (
		wg             sync.WaitGroup
		item           interface{}
		itemExits      bool
		itemIndex      int
		itemIndexExits bool
	)

	wg.Add(2)
	go func() {
		item, itemExits = c.items.Load(key)
		wg.Done()
	}()

	go func() {
		itemIndex, itemIndexExits = findIndexByItem(key, c.keys)
		wg.Done()
	}()
	wg.Wait()

	if itemExits && itemIndexExits {
		c.rwMutex.Lock()
		c.items.Delete(key)
		c.linearCurrentSizes -= int64(unsafe.Sizeof(item))
		c.keys = removeItemByIndex(c.keys, itemIndex) //Update keys slice after remove that key from items map
		c.rwMutex.Unlock()

		return item, nil
	}

	return nil, nil
}

// Read method return the item by key from linear without remove it
func (c *Client) Read(key string) (interface{}, error) {
	if c.IsEmpty() {
		return nil, errors.New("linear is empty")
	}

	item, exits := c.items.Load(key)
	if !exits {
		return nil, nil
	}

	return item, nil
}

// Update reassign value to the key
func (c *Client) Update(key string, value interface{}) error {
	newItemSize := int64(unsafe.Sizeof(key)) + int64(unsafe.Sizeof(value))
	if newItemSize > c.linearSizes || c.IsEmpty() {
		return errors.New("linear is empty or not enough space")
	}

	c.rwMutex.Lock()
	currentSize, exits := c.IsExits(key)
	if !exits {
		c.rwMutex.Unlock()
		return errors.New("key does not exit")
	}
	c.items.Store(key, value)
	c.linearCurrentSizes += currentSize - newItemSize
	c.rwMutex.Unlock()

	return nil
}

// Range the LinearClient
func (c *Client) Range(fn func(key, value interface{}) bool) {
	c.items.Range(fn)
}

// IsExits check key exits or not
func (c *Client) IsExits(key string) (int64, bool) {
	value, exits := c.items.Load(key)
	if !exits {
		return 0, false
	}

	return int64(unsafe.Sizeof(key)) + int64(unsafe.Sizeof(value)), true
}

// IsEmpty check linear size
func (c *Client) IsEmpty() bool {
	return len(c.keys) == 0
}

// GetItems return the map contain items
func (c *Client) GetItems() sync.Map {
	return c.items
}

// Getkeys return the list of key
func (c *Client) Getkeys() []string {
	return c.keys
}

// GetNumberOfKeys return the number of keys
func (c *Client) GetNumberOfKeys() int {
	return len(c.keys)
}

// GetLinearSizes return the linear size
func (c *Client) GetLinearSizes() int64 {
	return c.linearSizes
}

// SetLinearSizes change the linear size with new value
func (c *Client) SetLinearSizes(linearSizes int64) {
	c.rwMutex.RLock()
	c.linearSizes = linearSizes
	c.rwMutex.RUnlock()
}

// GetLinearCurrentSize return the current linear size
func (c *Client) GetLinearCurrentSize() int64 {
	return c.linearCurrentSizes
}

// removeItemByIndex remove item out of []string by index but maintains order, and return the new one
// Source: https://yourbasic.org/golang/delete-element-slice/
func removeItemByIndex(s []string, idx int) []string {
	copy(s[idx:], s[idx+1:]) // Shift s[idx+1:] left one index.
	s[len(s)-1] = ""         // Erase last element (write zero value).
	return s[:len(s)-1]      // Truncate s.
}

// findIndexByItem return index belong to the key
// Source: https://stackoverflow.com/questions/46745043/performance-of-for-range-in-go
func findIndexByItem(keyName string, items []string) (int, bool) {
	for index := range items {
		if keyName == items[index] {
			return index, true
		}
	}

	return -1, false
}
