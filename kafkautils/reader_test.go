package kafkautils

import (
	"context"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/kafkatest"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestConsumeReaders(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	topic := kafkatest.Topic(t)
	readerCfg := kafka.ReaderConfig{
		Brokers:  []string{kafkatest.GetBroker()},
		Topic:    topic,
		GroupID:  "test",
		MinBytes: 1,
		MaxBytes: 1 << 20,
	}
	c := func(ctx context.Context, r FetchCommitter) error {
		cancel()
		return nil
	}
	errFunc := func(ctx context.Context, err error) {
		testutils.ErrorErr(t, err)
	}
	ConsumeReaders(ctx, readerCfg, 1, c, errFunc)
}

func TestConsumeReadersError(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	topic := kafkatest.Topic(t)
	readerCfg := kafka.ReaderConfig{
		Brokers:  []string{kafkatest.GetBroker()},
		Topic:    topic,
		GroupID:  "test",
		MinBytes: 1,
		MaxBytes: 1 << 20,
	}
	c := func(ctx context.Context, r FetchCommitter) error {
		return errors.New("error")
	}
	errFunc := func(ctx context.Context, err error) {
		cancel()
	}
	ConsumeReaders(ctx, readerCfg, 1, c, errFunc)
}

type testFetchCommitter struct {
	fetch  Fetcher
	commit func(context.Context, ...kafka.Message) error
}

func (r *testFetchCommitter) FetchMessage(ctx context.Context) (kafka.Message, error) {
	return r.fetch(ctx)
}

func (r *testFetchCommitter) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	return r.commit(ctx, msgs...)
}
