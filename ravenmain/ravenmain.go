// Package ravenmain provides Raven/Sentry related utilities for a main package.
package ravenmain

import (
	raven "github.com/getsentry/raven-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/envutils"
	"github.com/siddhant2408/golang-libraries/errors"
)

// Init initialize Raven/Sentry.
func Init(dsn string, release string, env envutils.Env) (closeutils.F, error) {
	if env == envutils.Testing || env == envutils.Development {
		return func() {}, nil
	}
	err := raven.SetDSN(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "set DSN")
	}
	raven.SetRelease(release)
	raven.SetEnvironment(string(env))
	cl := func() {
		raven.Wait()
	}
	return cl, nil
}
