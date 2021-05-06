package amqputils_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/siddhant2408/golang-libraries/amqptest"
	"github.com/siddhant2408/golang-libraries/amqputils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/streadway/amqp"
)

func TestConsumer(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var pCalled testutils.CallCounter
	p := func(context.Context, amqp.Delivery) error {
		pCalled.Call()
		return nil
	}
	c := &amqputils.Consumer{
		Processor: p,
	}
	var aaAckCalled testutils.CallCounter
	aaAck := func(tag uint64, multiple bool) error {
		aaAckCalled.Call()
		cancel()
		return nil
	}
	aa := &testAMQPAcknowledger{
		testAMQPAcknowledgerAck: aaAck,
	}
	dlv := amqp.Delivery{
		Acknowledger: aa,
		Exchange:     "exchange",
		RoutingKey:   "routing_key",
		Headers: amqp.Table{
			"foo": "bar",
		},
	}
	ch := make(chan amqp.Delivery, 1)
	ch <- dlv
	err := c.Consume(ctx, ch)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	pCalled.AssertCalled(t)
	aaAckCalled.AssertCalled(t)
}

func TestConsumerErrorChannelClosed(t *testing.T) {
	ctx := context.Background()
	ch := make(chan amqp.Delivery)
	close(ch)
	c := &amqputils.Consumer{}
	err := c.Consume(ctx, ch)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestConsumerErrorProcessor(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	p := func(context.Context, amqp.Delivery) error {
		return errors.New("error")
	}
	var eCalled testutils.CallCounter
	e := func(ctx context.Context, err error) {
		eCalled.Call()
	}
	c := &amqputils.Consumer{
		Processor: p,
		Error:     e,
	}
	aaNack := func(tag uint64, multiple bool, requeue bool) error {
		cancel()
		return nil
	}
	aa := &testAMQPAcknowledger{
		testAMQPAcknowledgerNack: aaNack,
	}
	dlv := amqp.Delivery{
		Acknowledger: aa,
		Exchange:     "exchange",
		RoutingKey:   "routing_key",
		Headers: amqp.Table{
			"foo": "bar",
		},
	}
	ch := make(chan amqp.Delivery, 1)
	ch <- dlv
	err := c.Consume(ctx, ch)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	eCalled.AssertCalled(t)
}

func TestConsumerErrorProcessorWithAcknowledger(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	p := func(context.Context, amqp.Delivery) error {
		err := errors.New("error")
		err = amqputils.ErrorWithAcknowledger(err, amqputils.NackDiscard)
		return err
	}
	e := func(ctx context.Context, err error) {}
	c := &amqputils.Consumer{
		Processor: p,
		Error:     e,
	}
	aaNack := func(tag uint64, multiple bool, requeue bool) error {
		if requeue {
			t.Fatal("requeue")
		}
		cancel()
		return nil
	}
	aa := &testAMQPAcknowledger{
		testAMQPAcknowledgerNack: aaNack,
	}
	dlv := amqp.Delivery{
		Acknowledger: aa,
	}
	ch := make(chan amqp.Delivery, 1)
	ch <- dlv
	err := c.Consume(ctx, ch)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestConsumerErrorProcessorNotTemporary(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	p := func(context.Context, amqp.Delivery) error {
		err := errors.New("error")
		err = errors.WithTemporary(err, false)
		return err
	}
	e := func(ctx context.Context, err error) {}
	c := &amqputils.Consumer{
		Processor: p,
		Error:     e,
	}
	aaNack := func(tag uint64, multiple bool, requeue bool) error {
		if requeue {
			t.Fatal("requeue")
		}
		cancel()
		return nil
	}
	aa := &testAMQPAcknowledger{
		testAMQPAcknowledgerNack: aaNack,
	}
	dlv := amqp.Delivery{
		Acknowledger: aa,
	}
	ch := make(chan amqp.Delivery, 1)
	ch <- dlv
	err := c.Consume(ctx, ch)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestConsumerErrorAcknowledge(t *testing.T) {
	ctx := context.Background()
	p := func(context.Context, amqp.Delivery) error {
		return nil
	}
	c := &amqputils.Consumer{
		Processor: p,
	}
	aaAck := func(tag uint64, multiple bool) error {
		return errors.New("error")
	}
	aa := &testAMQPAcknowledger{
		testAMQPAcknowledgerAck: aaAck,
	}
	dlv := amqp.Delivery{
		Acknowledger: aa,
	}
	ch := make(chan amqp.Delivery, 1)
	ch <- dlv
	err := c.Consume(ctx, ch)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestRunConsumer(t *testing.T) {
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
	amqputils.RunConsumer(ctx, cg, tp, queue, p, 1, errFunc)
}

func TestRunConsumers(t *testing.T) {
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
	amqputils.RunConsumers(ctx, cg, tp, queue, p, 1, 1, errFunc)
}

func TestGetOriginalPublish(t *testing.T) {
	for _, tc := range []struct {
		name               string
		msg                amqp.Delivery
		expectedExchange   string
		expectedRoutingKey string
	}{
		{
			name: "Normal",
			msg: amqp.Delivery{
				Exchange:   "a",
				RoutingKey: "b",
			},
			expectedExchange:   "a",
			expectedRoutingKey: "b",
		},
		{
			name: "Death",
			msg: amqp.Delivery{
				Exchange:   "a",
				RoutingKey: "b",
				Headers: amqp.Table{
					"x-death": []interface{}{
						amqp.Table{
							"exchange":     "c",
							"routing-keys": []interface{}{"d"},
						},
						amqp.Table{
							"exchange":     "e",
							"routing-keys": []interface{}{"f"},
						},
					},
				},
			},
			expectedExchange:   "e",
			expectedRoutingKey: "f",
		},
		{
			name: "DeathsEmpty",
			msg: amqp.Delivery{
				Exchange:   "a",
				RoutingKey: "b",
				Headers: amqp.Table{
					"x-death": []interface{}{},
				},
			},
			expectedExchange:   "a",
			expectedRoutingKey: "b",
		},
		{
			name: "DeathInvalidType",
			msg: amqp.Delivery{
				Exchange:   "a",
				RoutingKey: "b",
				Headers: amqp.Table{
					"x-death": []interface{}{
						"invalid",
					},
				},
			},
			expectedExchange:   "a",
			expectedRoutingKey: "b",
		},
		{
			name: "DeathExchangeInvalidType",
			msg: amqp.Delivery{
				Exchange:   "a",
				RoutingKey: "b",
				Headers: amqp.Table{
					"x-death": []interface{}{
						amqp.Table{
							"exchange":     123,
							"routing-keys": []interface{}{"d"},
						},
					},
				},
			},
			expectedExchange:   "a",
			expectedRoutingKey: "b",
		},
		{
			name: "DeathRoutingKeysInvalidType",
			msg: amqp.Delivery{
				Exchange:   "a",
				RoutingKey: "b",
				Headers: amqp.Table{
					"x-death": []interface{}{
						amqp.Table{
							"exchange":     "c",
							"routing-keys": "invalid",
						},
					},
				},
			},
			expectedExchange:   "a",
			expectedRoutingKey: "b",
		},
		{
			name: "DeathRoutingKeysEmpty",
			msg: amqp.Delivery{
				Exchange:   "a",
				RoutingKey: "b",
				Headers: amqp.Table{
					"x-death": []interface{}{
						amqp.Table{
							"exchange":     "c",
							"routing-keys": []interface{}{},
						},
					},
				},
			},
			expectedExchange:   "a",
			expectedRoutingKey: "b",
		},
		{
			name: "DeathRoutingKeyInvalidType",
			msg: amqp.Delivery{
				Exchange:   "a",
				RoutingKey: "b",
				Headers: amqp.Table{
					"x-death": []interface{}{
						amqp.Table{
							"exchange":     "c",
							"routing-keys": []interface{}{123},
						},
					},
				},
			},
			expectedExchange:   "a",
			expectedRoutingKey: "b",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			exchange, routingKey := amqputils.GetOriginalPublish(tc.msg)
			if exchange != tc.expectedExchange {
				t.Fatalf("unexpected exchange: got %q, want %q", exchange, tc.expectedExchange)
			}
			if routingKey != tc.expectedRoutingKey {
				t.Fatalf("unexpected routing key: got %q, want %q", routingKey, tc.expectedRoutingKey)
			}
		})
	}
}

func newTestQueue(ctx context.Context, tb testing.TB, cg amqputils.ChannelGetter) string {
	tb.Helper()
	chn, err := cg(ctx)
	if err != nil {
		testutils.FatalErr(tb, err)
	}
	defer chn.Close() //nolint:errcheck
	name := fmt.Sprintf("test_%d", rand.Int63())
	_, err = chn.QueueDeclare(name, false, true, false, false, nil)
	if err != nil {
		testutils.FatalErr(tb, err)
	}
	return name
}
