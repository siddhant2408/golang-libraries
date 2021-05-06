package kafkautils

import (
	"context"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestBatchConsumer(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	a := func(context.Context, Fetcher) ([]kafka.Message, error) {
		msgs := make([]kafka.Message, 10)
		return msgs, nil
	}
	var prCalled testutils.CallCounter
	pr := func(ctx context.Context, msgs []kafka.Message) error {
		prCalled.Call()
		return nil
	}
	r := &testFetchCommitter{
		commit: func(context.Context, ...kafka.Message) error {
			cancel()
			return nil
		},
	}
	c := &BatchConsumer{
		Accumulator: a,
		Processor:   pr,
	}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	prCalled.AssertCalled(t)
}

func TestBatchConsumerAccumulatorContextDone(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	a := func(context.Context, Fetcher) ([]kafka.Message, error) {
		cancel()
		return nil, nil
	}
	r := &testFetchCommitter{}
	c := &BatchConsumer{
		Accumulator: a,
	}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestBatchConsumerAccumulatorNoMessage(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var calls int
	a := func(context.Context, Fetcher) ([]kafka.Message, error) {
		calls++
		if calls >= 2 {
			cancel()
		}
		return nil, nil
	}
	r := &testFetchCommitter{}
	c := &BatchConsumer{
		Accumulator: a,
	}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestBatchConsumerErrorAccumulator(t *testing.T) {
	ctx := context.Background()
	a := func(context.Context, Fetcher) ([]kafka.Message, error) {
		return nil, errors.New("error")
	}
	r := &testFetchCommitter{}
	c := &BatchConsumer{
		Accumulator: a,
	}
	err := c.Consume(ctx, r)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestBatchConsumerErrorProcessor(t *testing.T) {
	ctx := context.Background()
	a := func(context.Context, Fetcher) ([]kafka.Message, error) {
		msgs := make([]kafka.Message, 10)
		return msgs, nil
	}
	pr := func(ctx context.Context, msgs []kafka.Message) error {
		return errors.New("error")
	}
	r := &testFetchCommitter{}
	c := &BatchConsumer{
		Accumulator: a,
		Processor:   pr,
	}
	err := c.Consume(ctx, r)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestBatchConsumerErrorCommit(t *testing.T) {
	ctx := context.Background()
	a := func(context.Context, Fetcher) ([]kafka.Message, error) {
		msgs := make([]kafka.Message, 10)
		return msgs, nil
	}
	pr := func(ctx context.Context, msgs []kafka.Message) error {
		return nil
	}
	r := &testFetchCommitter{
		commit: func(context.Context, ...kafka.Message) error {
			return errors.New("error")
		},
	}
	c := &BatchConsumer{
		Accumulator: a,
		Processor:   pr,
	}
	err := c.Consume(ctx, r)
	if err == nil {
		t.Fatal("no error")
	}
}
