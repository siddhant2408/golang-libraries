// Package mongotest provides testing helpers for MongoDB.
//
// If MongoDB is not available, the test is skipped.
// It can be controlled with the MONGOTEST_UNAVAILABLE_SKIP environment variable.
package mongotest

import (
	"context"
	"fmt"
	"os"
	"reflect" //nolint:depguard // Required for the mock cursor.
	"sync"
	"testing"
	"time"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/mongo"
	"github.com/siddhant2408/golang-libraries/mongoutils"
	"github.com/siddhant2408/golang-libraries/testutils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	uriDefault = "mongodb://localhost:27017"
	uriEnvVar  = "MONGOTEST_URI"

	connectTimeoutDefault = 500 * time.Millisecond
	connectTimeoutEnvVar  = "MONGOTEST_CONNECT_TIMEOUT"

	unavailableSkipEnvVar = "MONGOTEST_UNAVAILABLE_SKIP"
)

// GetURI returns the URL for the local test instance.
// It can be overridden with the MONGOTEST_URI environment variable.
func GetURI() string {
	u, ok := os.LookupEnv(uriEnvVar)
	if ok {
		return u
	}
	return uriDefault
}

// Connect connects to the test instance and returns a client.
//
// It registers a cleanup function that closes the client at the end of the test.
func Connect(tb testing.TB) *mongo.Client {
	tb.Helper()
	checkAvailable(tb)
	ctx := context.Background()
	connectTimeout := getConnectTimeout(tb)
	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()
	opts := options.Client().ApplyURI(GetURI())
	clt, err := mongo.Connect(ctx, opts)
	if err != nil {
		setNotAvailable(tb, err)
	}
	tb.Cleanup(func() {
		_ = mongoutils.ForceDisconnect(clt)
	})
	err = clt.Ping(ctx, nil)
	if err != nil {
		setNotAvailable(tb, err)
	}
	return clt
}

func getConnectTimeout(tb testing.TB) time.Duration {
	tb.Helper()
	s, ok := os.LookupEnv(connectTimeoutEnvVar)
	if !ok {
		return connectTimeoutDefault
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		err = errors.Wrapf(err, "parse %q environment variable", connectTimeoutEnvVar)
		testutils.FatalErr(tb, err)
	}
	return d
}

var (
	muAvailable sync.Mutex
	available   = true
)

func checkAvailable(tb testing.TB) {
	tb.Helper()
	muAvailable.Lock()
	defer muAvailable.Unlock()
	if !available {
		err := errors.New("MongoDB is not available")
		testutils.HandleUnavailable(tb, unavailableSkipEnvVar, err)
	}
}

func setNotAvailable(tb testing.TB, err error) {
	tb.Helper()
	muAvailable.Lock()
	defer muAvailable.Unlock()
	available = false
	err = errors.Wrap(err, "MongoDB is not available")
	testutils.HandleUnavailable(tb, unavailableSkipEnvVar, err)
}

// CheckAvailable checks that the test instance is available.
func CheckAvailable(tb testing.TB) {
	tb.Helper()
	_ = Connect(tb)
}

// GetDatabase returns a test database.
//
// It registers a cleanup function that deletes the database at the end of the test.
func GetDatabase(tb testing.TB) *mongo.Database {
	tb.Helper()
	clt := Connect(tb)
	name := fmt.Sprintf("test_%s", primitive.NewObjectID().Hex())
	db := clt.Database(name)
	tb.Cleanup(func() {
		ctx := context.Background()
		_ = db.Drop(ctx)
	})
	return db
}

// Cursor is a mock for MongoDB cursor.
type Cursor struct {
	nextIndex int
	Documents []interface{}
	Error     error
}

// Next implements Cursor.
func (c *Cursor) Next(ctx context.Context) bool {
	ok := c.nextIndex < len(c.Documents)
	if ok {
		c.nextIndex++
		return true
	}
	return false
}

// Decode implements Cursor.
func (c *Cursor) Decode(res interface{}) error {
	resv := reflect.ValueOf(res)
	if !resv.IsValid() {
		return errors.New("result is not valid")
	}
	if resv.Kind() != reflect.Ptr {
		return errors.Newf("result must be a pointer: %T", res)
	}
	resve := resv.Elem()
	if !resve.IsValid() {
		return errors.New("result element is not valid")
	}
	index := c.nextIndex - 1
	if index < 0 || index >= len(c.Documents) {
		return errors.Newf("index out of bounds: index %d, length %d", index, len(c.Documents))
	}
	elem := c.Documents[index]
	elemv := reflect.ValueOf(elem)
	if !elemv.IsValid() {
		return errors.Newf("element is not valid: index %d", index)
	}
	if !elemv.Type().AssignableTo(resve.Type()) {
		return errors.Newf("unassignable element type: index %d, got %v, want %v", index, elemv.Type(), resve.Type())
	}
	resve.Set(elemv)
	return nil
}

// All implements Cursor.
func (c *Cursor) All(ctx context.Context, res interface{}) error {
	resv := reflect.ValueOf(res)
	if !resv.IsValid() {
		return errors.New("result is not valid")
	}
	if resv.Kind() != reflect.Ptr {
		return errors.Newf("result must be a pointer: %T", res)
	}
	resve := resv.Elem()
	if !resve.IsValid() {
		return errors.New("result element is not valid")
	}
	if resve.Kind() != reflect.Slice {
		return errors.Newf("result element must be a slice: %T", res)
	}
	elemType := resve.Type().Elem()
	resve = resve.Slice(0, 0)
	for c.Next(ctx) {
		elemvp := reflect.New(elemType)
		elemp := elemvp.Interface()
		err := c.Decode(elemp)
		if err != nil {
			return errors.Wrap(err, "decode")
		}
		elemv := elemvp.Elem()
		resve = reflect.Append(resve, elemv)
	}
	resv.Elem().Set(resve)
	return c.Err()
}

// Err implements Cursor.
func (c *Cursor) Err() error {
	return c.Error
}

// Close implements Cursor.
func (c *Cursor) Close(ctx context.Context) error {
	// No error.
	return nil
}
