// Package httptestlocal provides an HTTP RoundTripper that restricts outgoing HTTP request to localhost.
package httptestlocal

import (
	"context"
	"net"
	"net/http"

	"github.com/siddhant2408/golang-libraries/errors"
)

type roundTripper struct {
	http.RoundTripper
}

// WrapRoundTripper wraps the given RoundTripper.
func WrapRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &roundTripper{
		RoundTripper: rt,
	}
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if !rt.isAllowed(req) {
		if req.Body != nil {
			_ = req.Body.Close() // Required by the RoundTripper specification.
		}
		return nil, errors.New("external HTTP request is not allowed in test")
	}
	return rt.RoundTripper.RoundTrip(req)
}

func (rt *roundTripper) isAllowed(req *http.Request) bool {
	return isContextAllowed(req.Context()) || rt.isLocal(req)
}

func (rt *roundTripper) isLocal(req *http.Request) bool {
	h := req.URL.Hostname()
	ips, err := net.LookupIP(h)
	if err != nil {
		// We delegate this error to the child RoundTripper.
		return true
	}
	for _, ip := range ips {
		if !ip.IsLoopback() {
			return false
		}
	}
	return true
}

// WrapDefaultTransport initializes http.DefaultTransport.
// It wraps it with WrapRoundTripper.
func WrapDefaultTransport() {
	http.DefaultTransport = WrapRoundTripper(http.DefaultTransport)
}

type allowContextKey struct{}

// Allow allows the given context to do external HTTP request.
func Allow(ctx context.Context) context.Context {
	return context.WithValue(ctx, allowContextKey{}, struct{}{})
}

func isContextAllowed(ctx context.Context) bool {
	return ctx.Value(allowContextKey{}) != nil
}
