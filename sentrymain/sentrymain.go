// Package sentrymain provides Sentry related utilities for a main package.
package sentrymain

import (
	"github.com/getsentry/sentry-go"
	"github.com/siddhant2408/golang-libraries/envutils"
	"github.com/siddhant2408/golang-libraries/errors"
)

// Init initializes Sentry.
func Init(dsn string, release string, env envutils.Env) error {
	if env == envutils.Testing || env == envutils.Development {
		return nil
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Release:     release,
		Environment: string(env),
	})
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
