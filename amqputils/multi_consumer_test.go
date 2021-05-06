package amqputils_test

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/amqptest"
	"github.com/siddhant2408/golang-libraries/amqputils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/streadway/amqp"
)

func TestMultiConsumer(t *testing.T) {
	ctx := context.Background()
	var cCalled testutils.CallCounter
	c := func(ctx context.Context, ch <-chan amqp.Delivery) error {
		cCalled.Call()
		<-ch
		return nil
	}
	mc := &amqputils.MultiConsumer{
		Count:   8,
		Consume: c,
	}
	ch := make(chan amqp.Delivery, mc.Count)
	for i := 0; i < mc.Count; i++ {
		ch <- amqp.Delivery{}
	}
	err := mc.MultiConsume(ctx, ch)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	cCalled.AssertCount(t, int64(mc.Count))
}

func TestMultiConsumerError(t *testing.T) {
	ctx := context.Background()
	c := func(ctx context.Context, ch <-chan amqp.Delivery) error {
		return errors.New("error")
	}
	mc := &amqputils.MultiConsumer{
		Count:   8,
		Consume: c,
	}
	ch := make(chan amqp.Delivery)
	err := mc.MultiConsume(ctx, ch)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestRunMultiConsumer(t *testing.T) {
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
	err = chn.Publish("", queue, false, false, amqp.Publishing{
		Body: []byte("test"),
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	tp := amqputils.Topology{
		Queues: []amqputils.QueueConfig{
			{
				Name:       queue,
				AutoDelete: true,
			},
		},
	}
	p := func(ctx context.Context, dlv amqp.Delivery) error {
		cancel()
		return nil
	}
	errFunc := func(ctx context.Context, err error) {
		testutils.FatalErr(t, err)
	}
	amqputils.RunMultiConsumer(ctx, cg, tp, queue, p, 1, 1, errFunc)
}
