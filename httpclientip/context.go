package httpclientip

import (
	"context"
	"net"

	"github.com/siddhant2408/golang-libraries/errors"
)

// SetToContext sets the client IP to a context.
//
// It should be used only for testing purpose.
func SetToContext(ctx context.Context, ip net.IP) context.Context {
	return setToContext(ctx, &simpleContextValue{
		ip: ip,
	})
}

func setToContext(ctx context.Context, v contextValue) context.Context {
	return context.WithValue(ctx, contextKey{}, v)
}

// GetFromContext returns the client IP from a Context.
//
// It is required to use the Handler from this package, otherwise it returns an error.
func GetFromContext(ctx context.Context) (net.IP, error) {
	v, ok := ctx.Value(contextKey{}).(contextValue)
	if !ok {
		return nil, errors.New("context key not defined")
	}
	return v.getClientIP()
}

type contextKey struct{}

type contextValue interface {
	getClientIP() (net.IP, error)
}

type simpleContextValue struct {
	ip net.IP
}

func (v *simpleContextValue) getClientIP() (net.IP, error) {
	return v.ip, nil
}
