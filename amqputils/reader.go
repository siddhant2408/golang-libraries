package amqputils

import (
	"context"

	"github.com/siddhant2408/golang-libraries/ctxutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/goroutine"
	"github.com/streadway/amqp"
)

// Reader represents an AMQP reader.
// It allows to read messages on a queue using a channel.
type Reader struct {
	Channel ChannelGetter
	Start   ReaderStart
	Consume func(context.Context, <-chan amqp.Delivery) error
}

// Read reads messages.
func (r *Reader) Read(ctx context.Context) error {
	chn, err := r.Channel(ctx)
	if err != nil {
		return errors.Wrap(err, "channel")
	}
	defer chn.Close() //nolint:errcheck
	ch, err := r.Start(ctx, chn)
	if err != nil {
		return errors.Wrap(err, "start")
	}
	err = r.Consume(ctx, ch)
	if err != nil {
		return errors.Wrap(err, "consume")
	}
	return nil
}

// ReaderStart starts to consume messages on an amqp.Channel.
type ReaderStart func(context.Context, *amqp.Channel) (<-chan amqp.Delivery, error)

// NewReaderStartQueue returns a new ReaderStart for a queue.
func NewReaderStartQueue(queue string) ReaderStart {
	return func(ctx context.Context, chn *amqp.Channel) (<-chan amqp.Delivery, error) {
		ch, err := chn.Consume(queue, "", false, false, false, false, nil)
		if err != nil {
			return nil, errors.Wrap(err, "consume")
		}
		return ch, nil
	}
}

// NewReaderStartConsumer returns a ReaderStart for a Consumer.
//
// It initializes the topology + the prefetch.
func NewReaderStartConsumer(tp Topology, queue string, prefetch int) ReaderStart {
	startQueue := NewReaderStartQueue(queue)
	return func(ctx context.Context, chn *amqp.Channel) (<-chan amqp.Delivery, error) {
		err := tp.Init(ctx, chn)
		if err != nil {
			return nil, errors.Wrap(err, "topology")
		}
		err = chn.Qos(prefetch, 0, false)
		if err != nil {
			return nil, errors.Wrap(err, "set QOS")
		}
		return startQueue(ctx, chn)
	}
}

// RunReader runs the reader in a loop.
// The loop exits when the context is canceled.
// If an error is returned, errFunc is called.
func RunReader(ctx context.Context, r *Reader, errFunc func(context.Context, error)) {
	for !ctxutils.IsDone(ctx) {
		err := r.Read(ctx)
		if err != nil {
			err = errors.Wrap(err, "AMQP reader")
			errFunc(ctx, err)
		}
	}
}

// RunReaders calls RunReader concurrently.
func RunReaders(ctx context.Context, r *Reader, count int, errFunc func(context.Context, error)) {
	goroutine.RunN(count, func() {
		RunReader(ctx, r, errFunc)
	})
}
