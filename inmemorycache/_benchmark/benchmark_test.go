package benchmark

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	bluele_gcache "github.com/bluele/gcache"
	goburrow_cache "github.com/goburrow/cache"
	groupcache_lru "github.com/golang/groupcache/lru"
	hashicorp_lru "github.com/hashicorp/golang-lru"
	karlseguin_ccache "github.com/karlseguin/ccache"
	"github.com/siddhant2408/golang-libraries/inmemorycache"
	vitess_cache "github.com/youtube/vitess/go/cache"
)

const (
	populateSize = 100
	maxSize      = 1000
	expiration   = 5 * time.Minute
	globalKey    = "foo"
	globalValue  = "bar"
)

var getResult interface{}

func Benchmark(b *testing.B) {
	for _, tcl := range []struct {
		name       string
		new        func() *cache
		hasMaxSize bool
	}{
		{
			name: "SIB",
			new:  newSIB,
		},
		{
			name:       "SIBMaxSize",
			new:        newSIBMaxSize,
			hasMaxSize: true,
		},
		{
			name: "SIBTTL",
			new:  newSIBTTL,
		},
		{
			name:       "SIBMaxSizeTTL",
			new:        newSIBMaxSizeTTL,
			hasMaxSize: true,
		},
		{
			name:       "BlueleGcache",
			new:        newBlueleGcache,
			hasMaxSize: true,
		},
		{
			name:       "GoburrowCache",
			new:        newGoburrowCache,
			hasMaxSize: true,
		},
		{
			name:       "KarlseguinCcache",
			new:        newKarlseguinCcache,
			hasMaxSize: true,
		},
		{
			name:       "HashicorpLRU",
			new:        newHashicorpLRU,
			hasMaxSize: true,
		},
		{
			name:       "HashicorpARC",
			new:        newHashicorpARC,
			hasMaxSize: true,
		},
		{
			name:       "GroupcacheLRU",
			new:        newGroupcacheLRU,
			hasMaxSize: true,
		},
		{
			name:       "VitessCache",
			new:        newVitessCache,
			hasMaxSize: true,
		},
	} {
		b.Run(tcl.name, func(b *testing.B) {
			for _, tcb := range []struct {
				name        string
				bench       func(*testing.B, *cache)
				needMaxSize bool
			}{
				{
					name:  "Set",
					bench: benchmarkSet,
				},
				{
					name:  "Get",
					bench: benchmarkGet,
				},
				{
					name:  "SetGet",
					bench: benchmarkSetGet,
				},
				{
					name:        "SetEvict",
					bench:       benchmarkSetEvict,
					needMaxSize: true,
				},
				{
					name:        "SetGetEvict",
					bench:       benchmarkSetGetEvict,
					needMaxSize: true,
				},
			} {
				if tcb.needMaxSize && !tcl.hasMaxSize {
					continue
				}
				b.Run(tcb.name, func(b *testing.B) {
					c := tcl.new()
					tcb.bench(b, c)
					b.StopTimer()
					if c.close != nil {
						c.close()
					}
				})
			}
		})
	}
}

func newSIB() *cache {
	c := inmemorycache.New(
		inmemorycache.Concurrent(),
	)
	return &cache{
		setKeyInterface: func(key interface{}, value interface{}) {
			c.Set(key, value)
		},
		getKeyInterface: c.Get,
	}
}

func newSIBMaxSize() *cache {
	c := inmemorycache.New(
		inmemorycache.Concurrent(),
		inmemorycache.MaxSize(maxSize),
	)
	return &cache{
		setKeyInterface: func(key interface{}, value interface{}) {
			c.Set(key, value)
		},
		getKeyInterface: c.Get,
	}
}

func newSIBTTL() *cache {
	c := inmemorycache.New(
		inmemorycache.Concurrent(),
		inmemorycache.TTL(expiration),
	)
	return &cache{
		setKeyInterface: func(key interface{}, value interface{}) {
			c.Set(key, value)
		},
		getKeyInterface: c.Get,
	}
}

func newSIBMaxSizeTTL() *cache {
	c := inmemorycache.New(
		inmemorycache.Concurrent(),
		inmemorycache.MaxSize(maxSize),
		inmemorycache.TTL(expiration),
	)
	return &cache{
		setKeyInterface: func(key interface{}, value interface{}) {
			c.Set(key, value)
		},
		getKeyInterface: c.Get,
	}
}

func newBlueleGcache() *cache {
	c := bluele_gcache.
		New(maxSize).
		Expiration(expiration).
		LRU().
		Build()
	return &cache{
		setKeyInterface: func(key interface{}, value interface{}) {
			_ = c.Set(key, value)
		},
		getKeyInterface: func(key interface{}) (value interface{}, found bool) {
			value, err := c.GetIFPresent(key)
			if err != nil {
				return nil, false
			}
			return value, true
		},
	}
}

func newGoburrowCache() *cache {
	c := goburrow_cache.New(
		goburrow_cache.WithMaximumSize(maxSize),
		//goburrow_cache.WithExpireAfterWrite(expiration),
		goburrow_cache.WithPolicy("lru"),
	)
	return &cache{
		setKeyInterface: func(key interface{}, value interface{}) {
			c.Put(key, value)
		},
		getKeyInterface: func(key interface{}) (value interface{}, found bool) {
			return c.GetIfPresent(key)
		},
		close: func() {
			_ = c.Close()
		},
	}
}

func newKarlseguinCcache() *cache {
	c := karlseguin_ccache.New(
		karlseguin_ccache.
			Configure().
			MaxSize(maxSize),
	)
	return &cache{
		setKeyString: func(key string, value interface{}) {
			c.Set(key, value, expiration)
		},
		getKeyString: func(key string) (value interface{}, found bool) {
			it := c.Get(key)
			if it == nil {
				return nil, false
			}
			return it.Value(), true
		},
		close: c.Stop,
	}
}

func newHashicorpLRU() *cache {
	c, err := hashicorp_lru.New(maxSize)
	if err != nil {
		panic(err)
	}
	return &cache{
		setKeyInterface: func(key interface{}, value interface{}) {
			c.Add(key, value)
		},
		getKeyInterface: c.Get,
	}
}

func newHashicorpARC() *cache {
	c, err := hashicorp_lru.NewARC(maxSize)
	if err != nil {
		panic(err)
	}
	return &cache{
		setKeyInterface: func(key interface{}, value interface{}) {
			c.Add(key, value)
		},
		getKeyInterface: c.Get,
	}
}

func newGroupcacheLRU() *cache {
	c := groupcache_lru.New(maxSize)
	var mu sync.Mutex
	return &cache{
		setKeyInterface: func(key interface{}, value interface{}) {
			mu.Lock()
			c.Add(key, value)
			mu.Unlock()
		},
		getKeyInterface: func(key interface{}) (value interface{}, found bool) {
			mu.Lock()
			value, found = c.Get(key)
			mu.Unlock()
			return value, found
		},
	}
}

func newVitessCache() *cache {
	c := vitess_cache.NewLRUCache(maxSize)
	return &cache{
		setKeyString: func(key string, value interface{}) {
			c.Set(key, &vitessValue{value: value})
		},
		getKeyString: func(key string) (value interface{}, found bool) {
			v, found := c.Get(key)
			if !found {
				return nil, false
			}
			return v.(*vitessValue).value, true
		},
	}
}

type vitessValue struct {
	value interface{}
}

func (v *vitessValue) Size() int {
	return 1
}

func benchmarkSet(b *testing.B, c *cache) {
	set := c.newSet(globalKey)
	benchmarkCommon(b, c)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			set(globalValue)
		}
	})
}

func benchmarkGet(b *testing.B, c *cache) {
	c.newSet(globalKey)(globalValue)
	get := c.newGet(globalKey)
	benchmarkCommon(b, c)
	b.RunParallel(func(pb *testing.PB) {
		var res interface{}
		for pb.Next() {
			v, found := get()
			if found {
				res = v
			}
		}
		getResult = res
	})
}

func benchmarkSetGet(b *testing.B, c *cache) {
	set := c.newSet(globalKey)
	get := c.newGet(globalKey)
	benchmarkCommon(b, c)
	b.RunParallel(func(pb *testing.PB) {
		var res interface{}
		for pb.Next() {
			set(globalValue)
			v, found := get()
			if found {
				res = v
			}
		}
		getResult = res
	})
}

func benchmarkSetEvict(b *testing.B, c *cache) {
	set := c.newSet(globalKey)
	benchmarkCommon(b, c)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			set(i)
			i++
		}
	})
}

func benchmarkSetGetEvict(b *testing.B, c *cache) {
	set := c.newSet(globalKey)
	get := c.newGet(globalKey)
	benchmarkCommon(b, c)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		var res interface{}
		for pb.Next() {
			set(i)
			v, found := get()
			if found {
				res = v
			}
			i++
		}
		getResult = res
	})
}

func benchmarkCommon(b *testing.B, c *cache) {
	populate(c)
	b.ResetTimer()
}

func populate(c *cache) {
	for i := 0; i < populateSize; i++ {
		c.newSet(strconv.Itoa(i))(fmt.Sprintf("test %d", i))
	}
}

type cache struct {
	setKeyInterface func(key interface{}, value interface{})
	getKeyInterface func(key interface{}) (value interface{}, found bool)
	setKeyString    func(key string, value interface{})
	getKeyString    func(key string) (value interface{}, found bool)
	close           func()
}

func (c *cache) isKeyInterface() bool {
	return c.setKeyInterface != nil && c.getKeyInterface != nil
}

func (c *cache) isKeyString() bool {
	return c.setKeyString != nil && c.getKeyString != nil
}

func (c *cache) newSet(key string) func(value interface{}) {
	if c.isKeyInterface() {
		var key interface{} = key
		return func(value interface{}) {
			c.setKeyInterface(key, value)
		}
	}
	if c.isKeyString() {
		return func(value interface{}) {
			c.setKeyString(key, value)
		}
	}
	panic("unsupported")
}

func (c *cache) newGet(key string) func() (value interface{}, found bool) {
	if c.isKeyInterface() {
		var key interface{} = key
		return func() (value interface{}, found bool) {
			return c.getKeyInterface(key)
		}
	}
	if c.isKeyString() {
		return func() (value interface{}, found bool) {
			return c.getKeyString(key)
		}
	}
	panic("unsupported")
}
