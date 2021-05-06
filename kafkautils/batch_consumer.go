package kafkautils

import (
	"context"
	"time"

	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/ctxutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

// BatchConsumer consumes batches of messages from a Reader.
//
// The steps are:
//  - fetch a batch of messages from the accumulator, without committing them
//  - call the processor
//  - commit the messages
//
// If the processor returns an error, the method Consume() returns.
type BatchConsumer struct {
	Accumulator func(context.Context, Fetcher) ([]kafka.Message, error)
	Processor   BatchConsumerProcessor
}

// Consume consumes batches of messages from a reader.
func (c *BatchConsumer) Consume(ctx context.Context, r FetchCommitter) error {
	for !ctxutils.IsDone(ctx) {
		err := c.consume(ctx, r)
		if err != nil {
			return errors.Wrap(err, "batch consumer")
		}
	}
	return nil
}

func (c *BatchConsumer) consume(ctx context.Context, r FetchCommitter) error {
	msgs, err := c.Accumulator(ctx, r.FetchMessage)
	if err != nil {
		return errors.Wrap(err, "accumulator")
	}
	if ctxutils.IsDone(ctx) {
		return nil
	}
	if len(msgs) == 0 {
		return nil
	}
	err = c.consumeMessages(r, msgs)
	if err != nil {
		return errors.Wrap(err, "messages")
	}
	return nil
}

func (c *BatchConsumer) consumeMessages(r FetchCommitter, msgs []kafka.Message) (err error) {
	ctx := context.Background() // Don't want to be interrupted.
	span, spanFinish := startTraceRootSpan(&ctx, "batch_consumer", &err)
	defer spanFinish()
	tracingutils.SetSpanType(span, tracingutils.SpanTypeMessageConsumer)
	opentracing_ext.SpanKindConsumer.Set(span)
	setTraceSpanTag(span, "messages.count", len(msgs))
	err = c.processMessages(ctx, msgs)
	if err != nil {
		return errors.Wrap(err, "process")
	}
	err = c.commitMessages(ctx, r, msgs)
	if err != nil {
		return errors.Wrap(err, "commit")
	}
	return nil
}

func (c *BatchConsumer) processMessages(ctx context.Context, msgs []kafka.Message) (err error) {
	_, spanFinish := startTraceChildSpan(&ctx, "batch_consumer.process", &err)
	defer spanFinish()
	err = c.Processor(ctx, msgs)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (c *BatchConsumer) commitMessages(ctx context.Context, r FetchCommitter, msgs []kafka.Message) (err error) {
	span, spanFinish := startTraceChildSpan(&ctx, "batch_consumer.commit", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, tracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.SpanTypeMessageConsumer)
	opentracing_ext.SpanKindConsumer.Set(span)
	err = r.CommitMessages(ctx, msgs...)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

// BatchConsumerProcessor processes a batch of messages from a BatchConsumer.
type BatchConsumerProcessor func(context.Context, []kafka.Message) error

// RunBatchConsumers runs batch consumers.
func RunBatchConsumers(ctx context.Context, readerCfg kafka.ReaderConfig, pr BatchConsumerProcessor, count int, size int, timeout time.Duration, errFunc func(context.Context, error)) {
	a := &Accumulator{
		Size:    size,
		Timeout: timeout,
	}
	c := &BatchConsumer{
		Accumulator: a.Accumulate,
		Processor:   pr,
	}
	ConsumeReaders(ctx, readerCfg, count, c.Consume, errFunc)
}
