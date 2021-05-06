package inmemorycache

import (
	"sync"
)

type lock struct {
	mu sync.Mutex
	Cache
}

func (c *lock) Set(key interface{}, value interface{}) {
	c.mu.Lock()
	c.Cache.Set(key, value)
	c.mu.Unlock()
}

func (c *lock) Get(key interface{}) (value interface{}, found bool) {
	c.mu.Lock()
	value, found = c.Cache.Get(key)
	c.mu.Unlock()
	return value, found
}

func (c *lock) Remove(key interface{}) {
	c.mu.Lock()
	c.Cache.Remove(key)
	c.mu.Unlock()
}
