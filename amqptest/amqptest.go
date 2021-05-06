// Package amqptest provides testing helpers for AMQP.
//
// If AMQP is not available, the test is skipped.
// It can be controlled with the AMQPTEST_UNAVAILABLE_SKIP environment variable.
package amqptest

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/siddhant2408/golang-libraries/amqputils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/streadway/amqp"
)

const (
	defaultURL            = "amqp://localhost:5672"
	urlEnvVar             = "AMQPTEST_URL"
	unavailableSkipEnvVar = "AMQPTEST_UNAVAILABLE_SKIP"
)

// GetURL returns the URL for the local test instance.
// It can be overridden with the AMQPTEST_URL environment variable.
func GetURL() string {
	u, ok := os.LookupEnv(urlEnvVar)
	if ok {
		return u
	}
	return defaultURL
}

// GetURLVHost calls GetURL and defines the vhost.
func GetURLVHost(tb testing.TB, vhost string) string {
	tb.Helper()
	u := GetURL()
	if vhost == "" || vhost == "/" {
		return u
	}
	pu, err := url.Parse(u)
	if err != nil {
		testutils.FatalErr(tb, errors.Wrap(err, "parse URL"))
	}
	pu.Path = vhost
	return pu.String()
}

// CheckAvailable checks that the local test instance is available.
func CheckAvailable(tb testing.TB) {
	tb.Helper()
	conn := NewConnection(tb, "")
	_ = conn.Close()
}

// NewConnection returns a new test Connection.
//
// It registers a cleanup function that closes the Connection at the end of the test.
func NewConnection(tb testing.TB, vhost string) *amqp.Connection {
	tb.Helper()
	ctx := context.Background()
	Vhost(tb, vhost)
	u := GetURLVHost(tb, vhost)
	conn, err := amqputils.Dial(ctx, u)
	if err != nil {
		err = errors.Wrapf(err, "AMQP is not available on %q", u)
		testutils.HandleUnavailable(tb, unavailableSkipEnvVar, err)
	}
	tb.Cleanup(func() {
		_ = conn.Close()
	})
	return conn
}

// NewConnectionManager returns a new test ConnectionManager.
//
// It registers a cleanup function that closes the ConnectionManager at the end of the test.
func NewConnectionManager(tb testing.TB, vhost string) *amqputils.ConnectionManager {
	tb.Helper()
	cm := &amqputils.ConnectionManager{
		Dial: func(ctx context.Context) (*amqp.Connection, error) {
			return NewConnection(tb, vhost), nil
		},
	}
	tb.Cleanup(func() {
		_ = cm.Close()
	})
	return cm
}

// Vhost ensures that a vhost exists and that the guest user has access to it.
func Vhost(tb testing.TB, name string) {
	tb.Helper()
	if name == "" || name == "/" {
		return
	}
	clt := newRabbitMQAPIClient(tb)
	// For some unknown reason, calls to the RabbitMQ API can hangs indefinitely.
	// So we retry several times.
	for i := 0; i < rabbitMQPutVhostAttempts; i++ {
		resp, err := clt.PutVhost(name, rabbithole.VhostSettings{})
		if err != nil {
			tb.Logf("create vhost %q attempt %d", name, i)
			testutils.LogErr(tb, err)
			continue
		}
		handleRabbitMQAPIResponse(tb, resp)
		resp, err = clt.UpdatePermissionsIn(name, rabbitMQAPIUsername, rabbithole.Permissions{
			Configure: ".*",
			Write:     ".*",
			Read:      ".*",
		})
		if err != nil {
			tb.Logf("update vhost %q permissions attempt %d", name, i)
			testutils.LogErr(tb, err)
			continue
		}
		handleRabbitMQAPIResponse(tb, resp)
		return
	}
	tb.Fatalf("failed to ensure vhost %q after %d attempts", name, rabbitMQPutVhostAttempts)
}

const (
	rabbitMQPutVhostAttempts = 5

	rabbitMQAPIURL      = "http://localhost:15672"
	rabbitMQAPIUsername = "guest"
	rabbitMQAPIPassword = "guest"
	rabbitmQAPITimeout  = 10 * time.Second
)

func newRabbitMQAPIClient(tb testing.TB) *rabbithole.Client {
	tb.Helper()
	c, err := rabbithole.NewClient(rabbitMQAPIURL, rabbitMQAPIUsername, rabbitMQAPIPassword)
	if err != nil {
		testutils.FatalErr(tb, err)
	}
	c.SetTimeout(rabbitmQAPITimeout)
	_, err = c.Overview()
	if err != nil {
		err = errors.Wrapf(err, "RabbitMQ management API is not available on %q", rabbitMQAPIURL)
		testutils.HandleUnavailable(tb, unavailableSkipEnvVar, err)
	}
	return c
}

func handleRabbitMQAPIResponse(tb testing.TB, resp *http.Response) {
	tb.Helper()
	defer resp.Body.Close() //nolint:errcheck
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		testutils.FatalErr(tb, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		tb.Fatalf("unexpected status code: got %d, want 2xx: %s", resp.StatusCode, body)
	}
}
