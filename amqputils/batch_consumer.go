package amqputils

import (
	"context"
	"time"

	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/siddhant2408/golang-libraries/ctxutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/timeutils"
	"github.com/siddhant2408/golang-libraries/tracingutils"
	"github.com/streadway/amqp"
)

// BatchConsumer consume messages in batch.
type BatchConsumer struct {
	Accumulator func(context.Context, <-chan amqp.Delivery) ([]amqp.Delivery, error)
	Processor   BatchConsumerProcessor
}

// Consume consumes messages.
func (c *BatchConsumer) Consume(ctx context.Context, ch <-chan amqp.Delivery) error {
	for !ctxutils.IsDone(ctx) {
		dlvs, err := c.Accumulator(ctx, ch)
		if err != nil {
			return errors.Wrap(err, "accumulator")
		}
		if len(dlvs) == 0 {
			continue
		}
		err = c.process(dlvs)
		if err != nil {
			return errors.Wrap(err, "processor")
		}
	}
	return nil
}

func (c *BatchConsumer) process(dlvs []amqp.Delivery) (err error) {
	ctx := context.Background()
	span, spanFinish := startTraceRootSpan(&ctx, "batch_consumer", &err)
	defer spanFinish()
	tracingutils.SetSpanType(span, tracingutils.SpanTypeMessageConsumer)
	opentracing_ext.SpanKindConsumer.Set(span)
	setTraceSpanTag(span, "deliveries.count", len(dlvs))
	return c.Processor(ctx, dlvs)
}

// BatchConsumerProcessor represents a processor for BatchConsumer.
type BatchConsumerProcessor func(context.Context, []amqp.Delivery) error

// RunBatchConsumers runs batch consumers on a queue.
//
// It accumulates `size` messages or wait for `delay`, then it processes them in batch.
// Each consumers has its own AMQP channel and prefetch.
// The prefetch is equal to `size * 2`, so it can continue to accumulate messages while processing.
//
// The argument `x-priority` is used, in order to deliver messages in priority to the same consumers.
// It allows to receives all the messages in the same consumer, and improve the probability to group them efficiently.
func RunBatchConsumers(ctx context.Context, cg ChannelGetter, tp Topology, queue string, pr BatchConsumerProcessor, count int, size int, delay time.Duration, errFunc func(context.Context, error)) {
	a := &Accumulator{
		Size:  size,
		Delay: delay,
	}
	c := &BatchConsumer{
		Accumulator: a.Accumulate,
		Processor:   pr,
	}
	start := func(ctx context.Context, chn *amqp.Channel) (<-chan amqp.Delivery, error) {
		err := tp.Init(ctx, chn)
		if err != nil {
			return nil, errors.Wrap(err, "topology")
		}
		// size*2 ensures that we can still receive messages on this consumer while processing the accumulated messages.
		err = chn.Qos(size*2, 0, false)
		if err != nil {
			return nil, errors.Wrap(err, "qos")
		}
		ch, err := chn.Consume(queue, "", false, false, false, false, amqp.Table{
			"x-priority": -timeutils.Now().UnixNano(), // The oldest consumer has the highest priority.
		})
		if err != nil {
			return nil, errors.Wrap(err, "consume")
		}
		return ch, nil
	}
	r := &Reader{
		Channel: cg,
		Start:   start,
		Consume: c.Consume,
	}
	RunReaders(ctx, r, count, errFunc)
}
