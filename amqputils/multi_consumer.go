package amqputils

import (
	"context"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/goroutine"
	"github.com/streadway/amqp"
)

// MultiConsumer prefetch multipliers.
const (
	MultiConsumerPrefetchMultiplierMinimum = 1
	MultiConsumerPrefetchMultiplierNormal  = 5
)

// MultiConsumer runs multiple consumers.
type MultiConsumer struct {
	Count   int
	Consume func(context.Context, <-chan amqp.Delivery) error
}

// MultiConsume consumes messages with multiple consumers.
func (mc *MultiConsumer) MultiConsume(ctx context.Context, ch <-chan amqp.Delivery) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errCh := make(chan error, 1)
	goroutine.RunN(mc.Count, func() {
		err := mc.Consume(ctx, ch)
		if err != nil {
			cancel()
			select {
			case errCh <- err:
			default:
			}
		}
	})
	select {
	case err := <-errCh:
		return errors.WithStack(err)
	default:
		return nil
	}
}

// RunMultiConsumer runs consumer on a queue.
//
// All consumers share the same AMQP channel and prefetch.
func RunMultiConsumer(ctx context.Context, cg ChannelGetter, tp Topology, queue string, pr ConsumerProcessor, count int, prefetch int, errFunc func(context.Context, error)) {
	c := &Consumer{
		Processor: pr,
		Error:     errFunc,
	}
	mc := &MultiConsumer{
		Count:   count,
		Consume: c.Consume,
	}
	start := NewReaderStartConsumer(tp, queue, prefetch)
	r := &Reader{
		Channel: cg,
		Start:   start,
		Consume: mc.MultiConsume,
	}
	RunReader(ctx, r, errFunc)
}
