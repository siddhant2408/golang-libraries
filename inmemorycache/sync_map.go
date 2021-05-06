package inmemorycache

import (
	"sync"
)

type syncMap struct {
	values sync.Map
}

func (c *syncMap) Set(key interface{}, value interface{}) {
	c.values.Store(key, value)
}

func (c *syncMap) Get(key interface{}) (value interface{}, found bool) {
	return c.values.Load(key)
}

func (c *syncMap) Remove(key interface{}) {
	c.values.Delete(key)
}
