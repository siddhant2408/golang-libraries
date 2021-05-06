package amqputils_test

import (
	"context"
	"sync"
	"testing"

	"github.com/siddhant2408/golang-libraries/amqptest"
	"github.com/siddhant2408/golang-libraries/amqputils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/streadway/amqp"
)

func TestReader(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	conn := amqptest.NewConnection(t, testVhost)
	r := &amqputils.Reader{
		Channel: amqputils.NewChannelGetterConnection(conn),
		Start: func(ctx context.Context, chn *amqp.Channel) (<-chan amqp.Delivery, error) {
			return nil, nil
		},
		Consume: func(ctx context.Context, chn <-chan amqp.Delivery) error {
			cancel()
			return nil
		},
	}
	err := r.Read(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestReaderErrorChannel(t *testing.T) {
	ctx := context.Background()
	r := &amqputils.Reader{
		Channel: func(ctx context.Context) (*amqp.Channel, error) {
			return nil, errors.New("error")
		},
	}
	err := r.Read(ctx)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestReaderErrorConsume(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	r := &amqputils.Reader{
		Channel: amqputils.NewChannelGetterConnection(conn),
		Start: func(ctx context.Context, chn *amqp.Channel) (<-chan amqp.Delivery, error) {
			return nil, errors.New("error")
		},
	}
	err := r.Read(ctx)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestReaderErrorProcess(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	r := &amqputils.Reader{
		Channel: amqputils.NewChannelGetterConnection(conn),
		Start: func(ctx context.Context, chn *amqp.Channel) (<-chan amqp.Delivery, error) {
			return nil, nil
		},
		Consume: func(ctx context.Context, chn <-chan amqp.Delivery) error {
			return errors.New("error")
		},
	}
	err := r.Read(ctx)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestNewReaderConsumeQueue(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	cg := amqputils.NewChannelGetterConnection(conn)
	queue := newTestQueue(ctx, t, cg)
	rs := amqputils.NewReaderStartQueue(queue)
	chn, err := cg(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	defer chn.Close() //nolint:errcheck
	ch, err := rs(ctx, chn)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if ch == nil {
		t.Fatal("nil")
	}
}

func TestRunReader(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	conn := amqptest.NewConnection(t, testVhost)
	r := &amqputils.Reader{
		Channel: amqputils.NewChannelGetterConnection(conn),
		Start: func(ctx context.Context, chn *amqp.Channel) (<-chan amqp.Delivery, error) {
			return nil, nil
		},
		Consume: func(ctx context.Context, chn <-chan amqp.Delivery) error {
			cancel()
			return nil
		},
	}
	amqputils.RunReader(ctx, r, nil)
}

func TestRunReaderError(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	conn := amqptest.NewConnection(t, testVhost)
	r := &amqputils.Reader{
		Channel: amqputils.NewChannelGetterConnection(conn),
		Start: func(ctx context.Context, chn *amqp.Channel) (<-chan amqp.Delivery, error) {
			return nil, nil
		},
		Consume: func(ctx context.Context, chn <-chan amqp.Delivery) error {
			return errors.New("error")
		},
	}
	amqputils.RunReader(ctx, r, func(ctx context.Context, err error) {
		cancel()
	})
}

func TestRunReaders(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	conn := amqptest.NewConnection(t, testVhost)
	var wg sync.WaitGroup
	count := 10
	wg.Add(count)
	r := &amqputils.Reader{
		Channel: amqputils.NewChannelGetterConnection(conn),
		Start: func(ctx context.Context, chn *amqp.Channel) (<-chan amqp.Delivery, error) {
			return nil, nil
		},
		Consume: func(ctx context.Context, chn <-chan amqp.Delivery) error {
			wg.Done()
			wg.Wait()
			cancel()
			return nil
		},
	}
	amqputils.RunReaders(ctx, r, count, nil)
}
