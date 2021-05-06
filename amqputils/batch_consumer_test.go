package amqputils_test

import (
	"context"
	"testing"
	"time"

	"github.com/siddhant2408/golang-libraries/amqptest"
	"github.com/siddhant2408/golang-libraries/amqputils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/streadway/amqp"
)

func TestBatchConsumer(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan amqp.Delivery)
	dlvs := make([]amqp.Delivery, 10)
	c := &amqputils.BatchConsumer{
		Accumulator: func(context.Context, <-chan amqp.Delivery) ([]amqp.Delivery, error) {
			return dlvs, nil
		},
		Processor: func(context.Context, []amqp.Delivery) error {
			cancel()
			return nil
		},
	}
	err := c.Consume(ctx, ch)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestBatchConsumerEmpty(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan amqp.Delivery)
	c := &amqputils.BatchConsumer{
		Accumulator: func(context.Context, <-chan amqp.Delivery) ([]amqp.Delivery, error) {
			cancel()
			return nil, nil
		},
	}
	err := c.Consume(ctx, ch)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestBatchConsumerErrorAccumulator(t *testing.T) {
	ctx := context.Background()
	ch := make(chan amqp.Delivery)
	c := &amqputils.BatchConsumer{
		Accumulator: func(context.Context, <-chan amqp.Delivery) ([]amqp.Delivery, error) {
			return nil, errors.New("error")
		},
	}
	err := c.Consume(ctx, ch)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestBatchConsumerErrorProcessor(t *testing.T) {
	ctx := context.Background()
	ch := make(chan amqp.Delivery)
	dlvs := make([]amqp.Delivery, 10)
	c := &amqputils.BatchConsumer{
		Accumulator: func(context.Context, <-chan amqp.Delivery) ([]amqp.Delivery, error) {
			return dlvs, nil
		},
		Processor: func(context.Context, []amqp.Delivery) error {
			return errors.New("error")
		},
	}
	err := c.Consume(ctx, ch)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestRunBatchConsumer(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	conn := amqptest.NewConnection(t, testVhost)
	cg := amqputils.NewChannelGetterConnection(conn)
	queue := newTestQueue(ctx, t, cg)
	chn, err := cg(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	defer chn.Close() //nolint:errcheck
	size := 10
	for i := 0; i < size; i++ {
		err = chn.Publish("", queue, false, false, amqp.Publishing{
			Body: []byte("test"),
		})
		if err != nil {
			testutils.FatalErr(t, err)
		}
	}
	tp := amqputils.Topology{
		Queues: []amqputils.QueueConfig{
			{
				Name:       queue,
				AutoDelete: true,
			},
		},
	}
	p := func(ctx context.Context, dlvs []amqp.Delivery) error {
		cancel()
		if len(dlvs) != size {
			t.Errorf("unexpected deliveries count: got %d, want %d", len(dlvs), size)
		}
		return nil
	}
	errFunc := func(ctx context.Context, err error) {
		testutils.FatalErr(t, err)
	}
	amqputils.RunBatchConsumers(ctx, cg, tp, queue, p, 1, size, 1*time.Second, errFunc)
}
