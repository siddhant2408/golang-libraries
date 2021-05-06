// Package loadutils provides loading related utilities
package loadutils

import (
	"context"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/singleflight"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

// Loader implements a loading mechanism with cache and singleflight.
type Loader struct {
	Cache interface {
		Set(ctx context.Context, key string, value interface{}) error
		Get(ctx context.Context, key string) (value interface{}, found bool, err error)
	}
	group singleflight.Group
}

// Load loads a value for a given key.
// If it finds it in the cache, it returns it.
// Otherwise, it runs a singleflight request, that loads the value, stores it in the cache (if the returned `cache` variable is true), then returns it.
func (l *Loader) Load(ctx context.Context, key string, f func(ctx context.Context) (value interface{}, cache bool, err error)) (value interface{}, err error) {
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "loader", &err)
	defer spanFinish()
	span.SetTag("loader.key", key)
	var found bool
	value, found, err = l.cacheGet(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "cache get")
	}
	if found {
		return value, nil
	}
	value, err = l.singleflight(ctx, key, f)
	if err != nil {
		err = errors.Wrap(err, "singleflight")
		return nil, err
	}
	return value, nil
}

func (l *Loader) cacheGet(ctx context.Context, key string) (value interface{}, found bool, err error) {
	_, spanFinish := tracingutils.StartChildSpan(&ctx, "loader.cache.get", &err)
	defer spanFinish()
	return l.Cache.Get(ctx, key)
}

func (l *Loader) singleflight(ctx context.Context, key string, f func(ctx context.Context) (value interface{}, cache bool, err error)) (value interface{}, err error) {
	_, spanFinish := tracingutils.StartChildSpan(&ctx, "loader.singleflight", &err)
	defer spanFinish()
	value, err, _ = l.group.Do(ctx, key, func(ctx context.Context) (value interface{}, err error) {
		value, cache, err := l.load(ctx, f)
		if err != nil {
			return nil, errors.Wrap(err, "load")
		}
		if cache {
			err = l.cacheSet(ctx, key, value)
			if err != nil {
				return nil, errors.Wrap(err, "cache set")
			}
		}
		return value, nil
	})
	return value, err
}

func (l *Loader) load(ctx context.Context, f func(ctx context.Context) (value interface{}, cache bool, err error)) (value interface{}, cache bool, err error) {
	_, spanFinish := tracingutils.StartChildSpan(&ctx, "loader.load", &err)
	defer spanFinish()
	return f(ctx)
}

func (l *Loader) cacheSet(ctx context.Context, key string, value interface{}) (err error) {
	_, spanFinish := tracingutils.StartChildSpan(&ctx, "loader.cache.set", &err)
	defer spanFinish()
	return l.Cache.Set(ctx, key, value)
}

// InMemoryCacheWrapper is a wrapper that allows to use inmemorycache.Cache in Loader.
type InMemoryCacheWrapper struct {
	Cache interface {
		Set(key interface{}, value interface{})
		Get(key interface{}) (value interface{}, found bool)
	}
}

// Set implements Loader.Cache.
func (w *InMemoryCacheWrapper) Set(ctx context.Context, key string, value interface{}) error {
	w.Cache.Set(key, value)
	return nil
}

// Get implements Loader.Cache.
func (w *InMemoryCacheWrapper) Get(ctx context.Context, key string) (value interface{}, found bool, err error) {
	value, found = w.Cache.Get(key)
	return value, found, nil
}
