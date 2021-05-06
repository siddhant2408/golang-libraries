// Package redistest provides testing helpers for Redis.
//
// If Redis is not available, the test is skipped.
// It can be controlled with the REDISTEST_UNAVAILABLE_SKIP environment variable.
package redistest

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
)

const (
	defaultAddr           = "localhost:6379"
	addrEnvVar            = "REDISTEST_ADDR"
	defaultDatabase       = 15
	databaseEnvVar        = "REDISTEST_DB"
	unavailableSkipEnvVar = "REDISTEST_UNAVAILABLE_SKIP"
)

// GetAddress returns the address for the local test instance.
// It can be overridden with the REDISTEST_ADDR environment variable.
func GetAddress() string {
	a, ok := os.LookupEnv(addrEnvVar)
	if ok {
		return a
	}
	return defaultAddr
}

// GetDatabase returns the database for the local test instance.
// It can be overridden with the REDISTEST_DB environment variable.
func GetDatabase(tb testing.TB) int {
	dbs, ok := os.LookupEnv(databaseEnvVar)
	if ok {
		db, err := strconv.Atoi(dbs)
		if err != nil {
			testutils.FatalErr(tb, errors.Wrap(err, "parse string to int"))
		}
		return db
	}
	return defaultDatabase
}

// NewClient returns a new test client.
//
// It registers a cleanup function that flushes the database and closes the client at the end of the test.
func NewClient(tb testing.TB) *redis.Client {
	tb.Helper()
	ctx := context.Background()
	addr := GetAddress()
	db := GetDatabase(tb)
	c := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})
	tb.Cleanup(func() {
		_ = c.Close()
	})
	err := c.Ping(ctx).Err()
	if err != nil {
		err = errors.Wrapf(err, "Redis not available on %q", addr)
		testutils.HandleUnavailable(tb, unavailableSkipEnvVar, err)
	}
	err = c.FlushDB(ctx).Err()
	if err != nil {
		err = errors.Wrap(err, "flush DB")
		testutils.FatalErr(tb, err)
	}
	tb.Cleanup(func() {
		_ = c.FlushDB(ctx)
	})
	return c
}

// CheckAvailable checks that the local test instance is available.
func CheckAvailable(tb testing.TB) {
	tb.Helper()
	_ = NewClient(tb)
}
