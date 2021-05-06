// Package httpsib provides HTTP related utilities.
package httpsib

import (
	"net/http"

	"github.com/siddhant2408/golang-libraries/httpclientip"
	"github.com/siddhant2408/golang-libraries/httpclientip/httpclientipsib"
	"github.com/siddhant2408/golang-libraries/httphandler"
)

// WrapHandler wraps a Handler and adds useful features.
//
//  - Catch panic and log
//  - /_ping route for healthcheck
//  - Limit concurrent requests (default = unlimited)
//  - Context cancel (default = disabled)
//  - Maximum request body size (default = 1MiB)
func WrapHandler(h http.Handler, opts ...Option) http.Handler {
	o := newOptions()
	for _, opt := range opts {
		opt(o)
	}
	h = &httpclientip.Handler{
		Handler: h,
		Getter:  httpclientipsib.Getter,
	}
	if o.requestBodyMaxBytes > 0 {
		h = &httphandler.RequestBodyMaxBytes{
			Handler:  h,
			MaxBytes: o.requestBodyMaxBytes,
		}
	}
	h = &httphandler.Ping{
		Handler: h,
	}
	if !o.contextCancel {
		h = &httphandler.NoContextCancel{
			Handler: h,
		}
	}
	if o.maxConcurrentLimit > 0 {
		h = &httphandler.MaxConcurrent{
			Handler:             h,
			Limit:               o.maxConcurrentLimit,
			LimitReachedHandler: o.maxConcurrentLimitReachedHandler,
		}
	}
	h = &httphandler.Panic{
		Handler: h,
	}
	return h
}

// Option represent an option.
type Option func(*options)

// MaxConcurrent returns an Option that limits the number of concurrent requests.
func MaxConcurrent(limit int, limitReachedHandler http.Handler) Option {
	return func(opts *options) {
		opts.maxConcurrentLimit = limit
		opts.maxConcurrentLimitReachedHandler = limitReachedHandler
	}
}

// RequestBodyMaxBytes returns an Option that defines the request body max bytes.
func RequestBodyMaxBytes(maxBytes int64) Option {
	return func(opts *options) {
		opts.requestBodyMaxBytes = maxBytes
	}
}

// ContextCancel returns an Option that controls context cancel (if the connection is closed).
func ContextCancel(enabled bool) Option {
	return func(opts *options) {
		opts.contextCancel = enabled
	}
}

type options struct {
	maxConcurrentLimit               int
	maxConcurrentLimitReachedHandler http.Handler
	requestBodyMaxBytes              int64
	contextCancel                    bool
}

func newOptions() *options {
	return &options{
		requestBodyMaxBytes: requestBodyMaxBytesDefault,
	}
}

const (
	requestBodyMaxBytesDefault = 1 << 20 // 1 MiB
)
