// Package redisutils provides redis related utilities.
//
// The clients created by this package have an application name and tracing.
package redisutils

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/redistracing"
)

// NewClient is a wrapper.
func NewClient(opts *redis.Options, appName string) *redis.Client {
	opts.OnConnect = wrapOnConnectSetName(opts.OnConnect, appName)
	clt := redis.NewClient(opts)
	clt.AddHook(&redistracing.Hook{
		Addrs: []string{opts.Addr},
		DB:    &opts.DB,
	})
	return clt
}

// NewFailoverClient is a wrapper.
func NewFailoverClient(opts *redis.FailoverOptions, appName string) *redis.Client {
	opts.OnConnect = wrapOnConnectSetName(opts.OnConnect, appName)
	clt := redis.NewFailoverClient(opts)
	clt.AddHook(&redistracing.Hook{
		Addrs: opts.SentinelAddrs,
		DB:    &opts.DB,
	})
	return clt
}

// NewClusterClient is a wrapper.
func NewClusterClient(opts *redis.ClusterOptions, appName string) *redis.ClusterClient {
	opts.OnConnect = wrapOnConnectSetName(opts.OnConnect, appName)
	clt := redis.NewClusterClient(opts)
	clt.AddHook(&redistracing.Hook{
		Addrs: opts.Addrs,
	})
	return clt
}

// NewFailoverClusterClient is a wrapper.
func NewFailoverClusterClient(opts *redis.FailoverOptions, appName string) *redis.ClusterClient {
	opts.OnConnect = wrapOnConnectSetName(opts.OnConnect, appName)
	clt := redis.NewFailoverClusterClient(opts)
	clt.AddHook(&redistracing.Hook{
		Addrs: opts.SentinelAddrs,
		DB:    &opts.DB,
	})
	return clt
}

// NewRing is a wrapper.
func NewRing(opts *redis.RingOptions, appName string) *redis.Ring {
	opts.OnConnect = wrapOnConnectSetName(opts.OnConnect, appName)
	clt := redis.NewRing(opts)
	addrs := make([]string, 0, len(opts.Addrs))
	for _, addr := range opts.Addrs {
		addrs = append(addrs, addr)
	}
	clt.AddHook(&redistracing.Hook{
		Addrs: addrs,
		DB:    &opts.DB,
	})
	return clt
}

// NewSentinelClient is a wrapper.
func NewSentinelClient(opts *redis.Options, appName string) *redis.SentinelClient {
	opts.OnConnect = wrapOnConnectSetName(opts.OnConnect, appName)
	clt := redis.NewSentinelClient(opts)
	clt.AddHook(&redistracing.Hook{
		Addrs: []string{opts.Addr},
		DB:    &opts.DB,
	})
	return clt
}

// NewUniversalClient is a wrapper.
func NewUniversalClient(opts *redis.UniversalOptions, appName string) redis.UniversalClient {
	opts.OnConnect = wrapOnConnectSetName(opts.OnConnect, appName)
	clt := redis.NewUniversalClient(opts)
	clt.AddHook(&redistracing.Hook{
		Addrs: opts.Addrs,
		DB:    &opts.DB,
	})
	return clt
}

type onConnect func(context.Context, *redis.Conn) error

func newOnConnectSetName(appName string) onConnect {
	return func(ctx context.Context, c *redis.Conn) error {
		err := c.ClientSetName(ctx, appName).Err()
		return errors.Wrap(err, "client set name")
	}
}

func wrapOnConnectSetName(oc onConnect, appName string) onConnect {
	sn := newOnConnectSetName(appName)
	if oc == nil {
		return sn
	}
	return func(ctx context.Context, c *redis.Conn) error {
		err := oc(ctx, c)
		if err != nil {
			return err
		}
		return sn(ctx, c)
	}
}
