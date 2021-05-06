package kafkautils

import (
	"context"
	"fmt"

	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/ctxutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

// Consumer consumes messages from a Reader.
//
// The steps are:
//  - fetch a message without committing it
//  - call the processor
//  - commit the message
//
// If the processor returns an error, the handler used is either (first):
//  - defined by ConsumerErrorWithHandler()
//  - ConsumerDiscard if the error is not temporary
//  - ConsumerRetry otherwise
type Consumer struct {
	// Processor processes messages.
	Processor ConsumerProcessor
	// Retry produces messages that must be retried.
	// Usually it should write to the same topic consumed by the consumer.
	Retry Producer
	// Discard produces messages that must be discarded.
	Discard Producer
	// Error is called if the processor returns an error.
	Error func(context.Context, error)
}

// Consume consumes messages from a reader.
func (c *Consumer) Consume(ctx context.Context, r FetchCommitter) error {
	for !ctxutils.IsDone(ctx) {
		err := c.consume(ctx, r)
		if err != nil {
			return errors.Wrap(err, "consumer")
		}
	}
	return nil
}

func (c *Consumer) consume(ctx context.Context, r FetchCommitter) error {
	msg, err := r.FetchMessage(ctx)
	if err != nil {
		if ctxutils.IsDone(ctx) {
			return nil
		}
		return errors.Wrap(err, "fetch")
	}
	err = c.consumeMessage(r, msg)
	if err != nil {
		return errors.Wrap(err, "message")
	}
	return nil
}

func (c *Consumer) consumeMessage(r FetchCommitter, msg kafka.Message) (err error) {
	ctx := context.Background() // Don't want to be interrupted.
	span, spanFinish := startTraceRootSpan(&ctx, "consumer", &err)
	defer spanFinish()
	tracingutils.SetSpanType(span, tracingutils.SpanTypeMessageConsumer)
	opentracing_ext.SpanKindConsumer.Set(span)
	setTraceSpanTagsMessageConsumer(span, msg)
	err = c.processMessage(ctx, msg)
	if err != nil {
		var h ConsumerErrorHandler
		h, err = c.getErrorHandler(err)
		if !errors.IsIgnored(err) {
			err = wrapErrorValueMessageConsumer(err, msg)
			tracingutils.SetSpanError(span, err)
			err = errors.Wrap(err, "Kafka consumer")
			c.Error(ctx, err)
		}
		err = h.Handle(ctx, c, msg, err)
		if err != nil {
			return errors.Wrap(err, "handle processor error")
		}
	}
	err = c.commitMessage(ctx, r, msg)
	if err != nil {
		return errors.Wrap(err, "commit")
	}
	return nil
}

func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message) (err error) {
	_, spanFinish := startTraceChildSpan(&ctx, "consumer.process", &err)
	defer spanFinish()
	err = c.Processor(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (c *Consumer) getErrorHandler(err error) (ConsumerErrorHandler, error) {
	h := GetConsumerErrorHandler(err)
	if h != nil {
		return h, err
	}
	if !errors.IsTemporary(err) {
		return ConsumerDiscard, errors.Wrap(err, "Kafka discard not temporary")
	}
	return ConsumerRetry, err
}

func (c *Consumer) commitMessage(ctx context.Context, r FetchCommitter, msg kafka.Message) (err error) {
	span, spanFinish := startTraceChildSpan(&ctx, "consumer.commit", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, tracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.SpanTypeMessageConsumer)
	opentracing_ext.SpanKindConsumer.Set(span)
	err = r.CommitMessages(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

// ConsumerProcessor processes a message from a Consumer.
type ConsumerProcessor func(context.Context, kafka.Message) error

// ConsumerErrorHandler handles Consumer error.
type ConsumerErrorHandler interface {
	Handle(context.Context, *Consumer, kafka.Message, error) error
	String() string
}

// ConsumerRetry produces the message to Consumer.RetryProducer.
//
// It should be used for temporary errors, that can be retried immediately and indefinitely.
const ConsumerRetry = consumerRetry("retry")

type consumerRetry string

func (h consumerRetry) Handle(ctx context.Context, c *Consumer, msg kafka.Message, cErr error) (err error) {
	_, spanFinish := startTraceChildSpan(&ctx, "consumer.error.retry", &err)
	defer spanFinish()
	msg = CopyMessage(msg)
	err = c.Retry(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "retry producer")
	}
	return nil
}

func (h consumerRetry) String() string {
	return string(h)
}

// ConsumerDiscard produces the message to Consumer.DiscardProducer, if defined.
//
// It should be used for invalid messages.
const ConsumerDiscard = consumerDiscard("discard")

type consumerDiscard string

func (h consumerDiscard) Handle(ctx context.Context, c *Consumer, msg kafka.Message, cErr error) (err error) {
	_, spanFinish := startTraceChildSpan(&ctx, "consumer.error.discard", &err)
	defer spanFinish()
	if c.Discard == nil {
		return nil
	}
	msg = CopyMessage(msg)
	err = c.Discard(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "discard producer")
	}
	return nil
}

func (h consumerDiscard) String() string {
	return string(h)
}

// ConsumerNoop does nothing.
const ConsumerNoop = consumerNoop("noop")

type consumerNoop string

func (h consumerNoop) Handle(ctx context.Context, c *Consumer, msg kafka.Message, cErr error) (err error) {
	_, spanFinish := startTraceChildSpan(&ctx, "consumer.error.noop", &err)
	spanFinish()
	return nil
}

func (h consumerNoop) String() string {
	return string(h)
}

// ConsumerErrorWithHandler wraps an error with a ConsumerErrorHandler.
func ConsumerErrorWithHandler(err error, h ConsumerErrorHandler) error {
	if err == nil {
		return nil
	}
	return &consumerHandlerError{
		err: err,
		h:   h,
	}
}

type consumerHandlerError struct {
	err error
	h   ConsumerErrorHandler
}

func (err *consumerHandlerError) KafkaConsumerErrorHandler() ConsumerErrorHandler {
	return err.h
}

func (err *consumerHandlerError) WriteErrorMessage(w errors.Writer, verbose bool) bool {
	_, _ = w.WriteString("Kafka ")
	_, _ = w.WriteString(err.h.String())
	return true
}

func (err *consumerHandlerError) Error() string                 { return errors.Error(err) }
func (err *consumerHandlerError) Format(s fmt.State, verb rune) { errors.Format(err, s, verb) }
func (err *consumerHandlerError) Unwrap() error                 { return err.err }

// GetConsumerErrorHandler returns the ConsumerErrorHandler associated to the error.
func GetConsumerErrorHandler(err error) ConsumerErrorHandler {
	var werr *consumerHandlerError
	ok := errors.As(err, &werr)
	if ok {
		return werr.KafkaConsumerErrorHandler()
	}
	return nil
}

// RunConsumers runs consumers.
func RunConsumers(ctx context.Context, readerCfg kafka.ReaderConfig, pr ConsumerProcessor, count int, retry Producer, discard Producer, errFunc func(context.Context, error)) {
	c := &Consumer{
		Processor: pr,
		Retry:     retry,
		Discard:   discard,
		Error:     errFunc,
	}
	ConsumeReaders(ctx, readerCfg, count, c.Consume, errFunc)
}
