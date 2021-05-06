package inmemorycache

import (
	"time"

	"github.com/siddhant2408/golang-libraries/timeutils"
)

type ttl struct {
	Cache
	ttl time.Duration
}

func (c *ttl) Set(key interface{}, value interface{}) {
	expiresNS := timeutils.Now().Add(c.ttl).UnixNano()
	ev := expiresValue{
		expiresNS: expiresNS,
		value:     value,
	}
	c.Cache.Set(key, ev)
}

func (c *ttl) Get(key interface{}) (value interface{}, found bool) {
	value, found = c.Cache.Get(key)
	if !found {
		return nil, false
	}
	// We're not checking the type assertion, because it shouldn't be possible that it fails.
	// If the type assertion fails, it means that the developer is doing something wrong with the internal structure.
	// So it's better that it panics, in order to notify the developer.
	ev := value.(expiresValue) //nolint:errcheck
	if ev.expiresNS < timeutils.Now().UnixNano() {
		c.Cache.Remove(key)
		return nil, false
	}
	return ev.value, true
}

type expiresValue struct {
	expiresNS int64
	value     interface{}
}
