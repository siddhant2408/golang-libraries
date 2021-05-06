package amqputils

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/siddhant2408/golang-libraries/ctxutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
	"github.com/streadway/amqp"
)

// Consumer consumes messages.
type Consumer struct {
	// Processor processes a single message.
	// It MUST NOT (n)ack the amqp.Delivery itself.
	// The Consumer will call the Acknowledger either (first):
	//  - Ack if error is nil
	//  - defined by ErrorWithAcknowledger()
	//  - NackDiscard if error is not temporary
	//  - NackRequeue
	Processor ConsumerProcessor
	// Error is called if the processor returns an error.
	Error func(context.Context, error)
}

// Consume consumes a channel of messages.
//
// It blocks until:
//  - The context is cancelled
//  - The channel is closed
//  - An internal error occurs
func (c *Consumer) Consume(ctx context.Context, ch <-chan amqp.Delivery) error {
	for !ctxutils.IsDone(ctx) {
		select {
		case dlv, ok := <-ch:
			if !ok {
				return errors.New("channel closed unexpectedly")
			}
			err := c.consumeDlv(dlv)
			if err != nil {
				return errors.Wrap(err, "delivery")
			}
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}

func (c *Consumer) consumeDlv(dlv amqp.Delivery) (err error) {
	ctx := context.Background()
	span, spanFinish := startTraceRootSpan(&ctx, "consumer", &err)
	defer spanFinish()
	c.updateTracingSpan(span, dlv)
	err = c.process(ctx, dlv)
	if err != nil && !errors.IsIgnored(err) {
		err = c.wrapProcessError(err, dlv)
		tracingutils.SetSpanError(span, err)
		err = errors.Wrap(err, "AMQP consumer")
		c.Error(ctx, err)
	}
	err = c.acknowledge(ctx, dlv, err)
	if err != nil {
		return errors.Wrap(err, "acknowledge")
	}
	return nil
}

func (c *Consumer) updateTracingSpan(span opentracing.Span, dlv amqp.Delivery) {
	tracingutils.SetSpanType(span, tracingutils.SpanTypeMessageConsumer)
	opentracing_ext.SpanKindConsumer.Set(span)
	if dlv.Exchange != "" {
		setTraceSpanTag(span, "exchange", dlv.Exchange)
	}
	if dlv.RoutingKey != "" {
		setTraceSpanTag(span, "routing_key", dlv.RoutingKey)
	}
	if len(dlv.Headers) > 0 {
		setTraceSpanTag(span, "headers", fmt.Sprint(dlv.Headers))
	}
	setTraceSpanTagBody(span, dlv.Body)
}

func (c *Consumer) process(ctx context.Context, dlv amqp.Delivery) (err error) {
	_, spanFinish := startTraceChildSpan(&ctx, "consumer.process", &err)
	defer spanFinish()
	// Don't allow the processor to acknowledge the message, it is not its role.
	dlv.Acknowledger = nil
	return c.Processor(ctx, dlv)
}

func (c *Consumer) wrapProcessError(err error, dlv amqp.Delivery) error {
	if dlv.Exchange != "" {
		err = wrapErrorValue(err, "exchange", dlv.Exchange)
	}
	if dlv.RoutingKey != "" {
		err = wrapErrorValue(err, "routing_key", dlv.RoutingKey)
	}
	if len(dlv.Headers) > 0 {
		err = wrapErrorValue(err, "headers", dlv.Headers)
	}
	err = wrapErrorValueBody(err, dlv.Body)
	if GetErrorAcknowledger(err) == nil && !errors.IsTemporary(err) {
		// If an error doesn't have an acknowledger and is NOT temporary, the message is discarded.
		err = ErrorWithAcknowledger(err, NackDiscard)
		err = errors.Wrap(err, "discard not temporary")
	}
	err = errors.Wrap(err, "process")
	return err
}

func (c *Consumer) acknowledge(ctx context.Context, dlv amqp.Delivery, myerr error) (err error) {
	span, spanFinish := startTraceChildSpan(&ctx, "consumer.acknowledge", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, tracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.SpanTypeMessageConsumer)
	opentracing_ext.SpanKindConsumer.Set(span)
	a := c.getAcknowledger(myerr)
	setTraceSpanTag(span, "acknowledger", a.String())
	return a.Acknowledge(dlv)
}

func (c *Consumer) getAcknowledger(err error) Acknowledger {
	if err == nil {
		return Ack
	}
	ack := GetErrorAcknowledger(err)
	if ack != nil {
		return ack
	}
	return NackRequeue
}

// ConsumerProcessor represents a processor for Consumer.
type ConsumerProcessor func(context.Context, amqp.Delivery) error

// RunConsumer runs a single consumer on a queue.
func RunConsumer(ctx context.Context, cg ChannelGetter, tp Topology, queue string, pr ConsumerProcessor, prefetch int, errFunc func(context.Context, error)) {
	c := &Consumer{
		Processor: pr,
		Error:     errFunc,
	}
	start := NewReaderStartConsumer(tp, queue, prefetch)
	r := &Reader{
		Channel: cg,
		Start:   start,
		Consume: c.Consume,
	}
	RunReader(ctx, r, errFunc)
}

// RunConsumers runs consumers on a queue.
//
// Each consumer has its own AMQP channel and prefetch.
func RunConsumers(ctx context.Context, cg ChannelGetter, tp Topology, queue string, pr ConsumerProcessor, count int, prefetch int, errFunc func(context.Context, error)) {
	c := &Consumer{
		Processor: pr,
		Error:     errFunc,
	}
	start := NewReaderStartConsumer(tp, queue, prefetch)
	r := &Reader{
		Channel: cg,
		Start:   start,
		Consume: c.Consume,
	}
	RunReaders(ctx, r, count, errFunc)
}

// GetOriginalPublish returns the original publish parameters for a message.
// It uses values from the header "x-death" if available.
func GetOriginalPublish(dlv amqp.Delivery) (exchange string, routingKey string) {
	exchange, routingKey, ok := getOriginalPublishFromDeath(dlv)
	if ok {
		return exchange, routingKey
	}
	return dlv.Exchange, dlv.RoutingKey
}

func getOriginalPublishFromDeath(dlv amqp.Delivery) (exchange string, routingKey string, ok bool) {
	deaths, ok := dlv.Headers["x-death"].([]interface{})
	if !ok {
		return "", "", false
	}
	if len(deaths) == 0 {
		return "", "", false
	}
	death, ok := deaths[len(deaths)-1].(amqp.Table)
	if !ok {
		return "", "", false
	}
	exchange, ok = death["exchange"].(string)
	if !ok {
		return "", "", false
	}
	rks, ok := death["routing-keys"].([]interface{})
	if !ok {
		return "", "", false
	}
	if len(rks) == 0 {
		return "", "", false
	}
	routingKey, ok = rks[0].(string)
	if !ok {
		return "", "", false
	}
	return exchange, routingKey, true
}
