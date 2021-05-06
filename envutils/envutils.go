// Package envutils provides environment related utilities.
package envutils

import (
	"flag"

	"github.com/siddhant2408/golang-libraries/errors"
)

const (
	// Testing is the testing environment.
	Testing = Env("testing")
	// Development is the development environment.
	Development = Env("development")
	// Staging is the staging environment.
	Staging = Env("staging")
	// Production is the production environment.
	Production = Env("production")
)

// Env represents an environment.
type Env string

func (e Env) String() string {
	return string(e)
}

// Set sets the string to the Env.
// It checks that the Env is valid.
func (e *Env) Set(s string) error {
	err := Check(Env(s))
	if err != nil {
		return err
	}
	*e = Env(s)
	return nil
}

// Check checks if an Env is valid.
func Check(env Env) error {
	switch env {
	case Testing, Development, Staging, Production:
		return nil
	default:
		return errors.Newf("unknown environment %q", env)
	}
}

// SetFlag sets the environment to the FlagSet.
func SetFlag(fs *flag.FlagSet, e *Env) {
	fs.Var(e, "environment", "Environment")
}
