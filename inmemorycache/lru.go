package inmemorycache

import (
	"container/list"
)

// This is inspired by https://github.com/golang/groupcache/blob/master/lru/lru.go .

type lru struct {
	maxSize int
	ll      *list.List
	cache   map[interface{}]*list.Element
}

type lruEntry struct {
	key   interface{}
	value interface{}
}

func newLRU(maxSize int) *lru {
	return &lru{
		maxSize: maxSize,
		ll:      list.New(),
		cache:   make(map[interface{}]*list.Element),
	}
}

func (c *lru) Set(key interface{}, value interface{}) {
	e, ok := c.cache[key]
	if ok {
		c.ll.MoveToFront(e)
		e.Value.(*lruEntry).value = value
		return
	}
	e = c.ll.PushFront(&lruEntry{key, value})
	c.cache[key] = e
	if c.ll.Len() > c.maxSize {
		c.removeOldest()
	}
}

func (c *lru) Get(key interface{}) (value interface{}, found bool) {
	e, ok := c.cache[key]
	if ok {
		c.ll.MoveToFront(e)
		return e.Value.(*lruEntry).value, true
	}
	return nil, false
}

func (c *lru) Remove(key interface{}) {
	e, ok := c.cache[key]
	if ok {
		c.removeElement(e)
	}
}

func (c *lru) removeOldest() {
	e := c.ll.Back()
	if e != nil {
		c.removeElement(e)
	}
}

func (c *lru) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*lruEntry) //nolint:errcheck
	delete(c.cache, kv.key)
}
