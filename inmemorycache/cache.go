package inmemorycache

import (
	"time"
)

// Cache represents a cache.
type Cache interface {
	// Set sets a value to the cache.
	Set(key interface{}, value interface{})
	// Get looks up a key's value from the cache.
	Get(key interface{}) (value interface{}, found bool)
	// Remove removes the provided key from the cache.
	Remove(key interface{})
}

// New returns a new Cache.
func New(opts ...Option) Cache {
	cfg := getConfig(opts...)
	if useSyncMap(cfg) {
		// If the cache is simply concurrent, then a simpler implementation can be used.
		return &syncMap{}
	}
	var c Cache
	if cfg.maxSize == 0 {
		c = newSimple()
	} else {
		c = newLRU(cfg.maxSize)
	}
	if cfg.ttl != 0 {
		c = &ttl{
			Cache: c,
			ttl:   cfg.ttl,
		}
	}
	// The `lock` implementation must be done last, because it synchronizes access if multiple calls to the underlying cache are done.
	// E.g. the `ttl` implementation removes a key after fetching it, if it is expired.
	if cfg.concurrent {
		c = &lock{
			Cache: c,
		}
	}
	return c
}

func useSyncMap(cfg *config) bool {
	return !cfg.disableSyncMap && cfg.maxSize == 0 && cfg.ttl == 0 && cfg.concurrent
}

type config struct {
	maxSize        int
	ttl            time.Duration
	concurrent     bool
	disableSyncMap bool
}

func getConfig(opts ...Option) *config {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// Option represents an option for New.
type Option func(*config)

// MaxSize returns an Option that configures the cache maximum size.
// By default the cache is not limited in size.
//
// The current implementation uses an LRU cache.
func MaxSize(maxSize int) Option {
	return func(cfg *config) {
		cfg.maxSize = maxSize
	}
}

// TTL returns an Option that configures the cache TTL.
// By default the cache doesn't have a TTL.
func TTL(ttl time.Duration) Option {
	return func(cfg *config) {
		cfg.ttl = ttl
	}
}

// Concurrent returns an Option that configures the cache to be concurrent safe.
// By default the cache is not concurrent safe.
func Concurrent() Option {
	return func(cfg *config) {
		cfg.concurrent = true
	}
}

func disableSyncMap() Option {
	return func(cfg *config) {
		cfg.disableSyncMap = true
	}
}
