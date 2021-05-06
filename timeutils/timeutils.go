// Package timeutils provides utilities to fake time.
package timeutils

import (
	"sync"
	"time"

	"github.com/siddhant2408/golang-libraries/errors"
)

// Now returns the current time (real by default).
func Now() time.Time {
	return now()
}

var now func() time.Time

func init() {
	InitReal()
}

// InitReal initializes the time to real time.
func InitReal() {
	now = time.Now
}

// InitFixed initializes the time to a fixed value.
func InitFixed() {
	SetFixed(time.Now())
}

// SetFixed sets the time to a fixed value.
func SetFixed(t time.Time) {
	now = func() time.Time {
		return t
	}
}

// Since is a replacement for `time.Since()`.
func Since(t time.Time) time.Duration {
	return Now().Sub(t)
}

// Until is a replacement for `time.Until()`.
func Until(t time.Time) time.Duration {
	return t.Sub(Now())
}

var locations sync.Map

// LoadLocation is a wrapper for time.LoadLocation.
// It caches the loaded locations, which is faster and decreases the memory allocations.
func LoadLocation(name string) (*time.Location, error) {
	v, ok := locations.Load(name)
	if ok {
		return v.(*time.Location), nil
	}
	l, err := time.LoadLocation(name)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	locations.Store(name, l)
	return l, nil
}
