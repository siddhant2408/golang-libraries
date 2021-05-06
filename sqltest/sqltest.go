// Package sqltest provides a helper to build SQL test packages
package sqltest

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/siddhant2408/golang-libraries/timeutils"
)

// Helper helps to build a SQL test package.
type Helper struct {
	Name                  string
	DSNDefault            string
	DSNEnvVar             string
	UnavailableSkipEnvVar string
	Open                  func(ctx context.Context, dsn string, dbName string) (*sql.DB, error)
	CreateQuery           func(dbName string) string
	DropQuery             func(dbName string) string
}

// GetDSN returns the DSN for the local test instance.
func (h *Helper) GetDSN() string {
	dsn, ok := os.LookupEnv(h.DSNEnvVar)
	if ok {
		return dsn
	}
	return h.DSNDefault
}

// GetDatabase returns a test database.
//
// It registers a cleanup function that deletes the database at the end of the test.
func (h *Helper) GetDatabase(tb testing.TB) *sql.DB {
	tb.Helper()
	dbName := fmt.Sprintf("test_%d_%d", timeutils.Now().UnixNano(), rand.Int63())
	h.createDatabase(tb, dbName)
	db := h.openCheck(tb, dbName)
	tb.Cleanup(func() {
		tb.Helper()
		_ = db.Close()
		h.dropDatabase(tb, dbName)
	})
	return db
}

// CheckAvailable checks that the local test instance is available.
func (h *Helper) CheckAvailable(tb testing.TB) {
	db := h.openCheck(tb, "")
	_ = db.Close()
}

func (h *Helper) openCheck(tb testing.TB, dbName string) *sql.DB {
	tb.Helper()
	ctx := context.Background()
	dsn := h.GetDSN()
	db, err := h.Open(ctx, dsn, dbName)
	if err != nil {
		err = errors.Wrap(err, "open")
		testutils.FatalErr(tb, err)
	}
	err = db.PingContext(ctx)
	if err != nil {
		err = errors.Wrapf(err, "%s is not available", h.Name)
		testutils.HandleUnavailable(tb, h.UnavailableSkipEnvVar, err)
	}
	return db
}

func (h *Helper) createDatabase(tb testing.TB, dbName string) {
	tb.Helper()
	ctx := context.Background()
	db := h.openCheck(tb, "")
	defer db.Close() //nolint:errcheck
	// Yes it's bad, but it's not possible to use placeholder in "CREATE DATABASE".
	_, err := db.ExecContext(ctx, h.CreateQuery(dbName))
	if err != nil {
		err = errors.Wrap(err, "create database")
		testutils.FatalErr(tb, err)
	}
}

func (h *Helper) dropDatabase(tb testing.TB, dbName string) {
	tb.Helper()
	ctx := context.Background()
	// Cancel the database drop after a short delay.
	// It may block indefinitely if there is a transaction locking the database.
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	db := h.openCheck(tb, "")
	defer db.Close() //nolint:errcheck
	// Yes it's bad, but it's not possible to use placeholder in "DROP DATABASE".
	_, err := db.ExecContext(ctx, h.DropQuery(dbName))
	if err != nil {
		err = errors.Wrap(err, "drop database")
		testutils.FatalErr(tb, err)
	}
}
