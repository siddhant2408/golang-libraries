package redislock

import (
	"github.com/siddhant2408/golang-libraries/errors"
)

const (
	// ErrNotObtained means that the lock could not be obtained.
	ErrNotObtained errorString = "not obtained"
	// ErrNotHeld means that the lock was not currently held.
	ErrNotHeld errorString = "not held"
)

type errorString string

func (err errorString) Error() string {
	return string(err)
}

func wrapErrorValue(err error, key string, val interface{}) error {
	return errors.WithValue(err, "redislock."+key, val)
}
