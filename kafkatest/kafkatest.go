// Package kafkatest provides Kafka test utilities.
//
// If Kafka is not available, the test is skipped.
// It can be controlled with the KAFKATEST_UNAVAILABLE_SKIP environment variable.
package kafkatest

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/siddhant2408/golang-libraries/timeutils"
)

const (
	defaultBroker         = "localhost:9092"
	brokerEnvVar          = "KAFKATEST_BROKER"
	unavailableSkipEnvVar = "KAFKATEST_UNAVAILABLE_SKIP"
)

// GetBroker returns the broker for the local test instance.
// It can be overridden with the KAFKATEST_BROKER environment variable.
func GetBroker() string {
	br, ok := os.LookupEnv(brokerEnvVar)
	if ok {
		return br
	}
	return defaultBroker
}

// CheckAvailable checks if the test broker is available.
func CheckAvailable(tb testing.TB) {
	tb.Helper()
	br := GetBroker()
	conn, err := kafka.Dial("tcp", br)
	if err != nil {
		err = errors.Wrapf(err, "Kafka is not available on %q", br)
		testutils.HandleUnavailable(tb, unavailableSkipEnvVar, err)
	}
	err = conn.Close()
	if err != nil {
		testutils.FatalErr(tb, err)
	}
}

// Topic returns a new test topic.
//
// It registers a cleanup function that deletes the topic at the end of the test.
func Topic(tb testing.TB) string {
	ctx := context.Background()
	tb.Helper()
	CheckAvailable(tb)
	c := &kafka.Client{
		Addr: kafka.TCP(GetBroker()),
	}
	name := fmt.Sprintf("test_%d_%d", timeutils.Now().UnixNano(), rand.Int63())
	cfg := kafka.TopicConfig{
		Topic:             name,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	_, err := c.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{cfg},
	})
	if err != nil {
		testutils.FatalErr(tb, err)
	}
	tb.Cleanup(func() {
		_, _ = c.DeleteTopics(ctx, &kafka.DeleteTopicsRequest{
			Topics: []string{name},
		})
	})
	return name
}
