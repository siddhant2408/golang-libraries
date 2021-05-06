package amqputils

import (
	"context"
	"fmt"
	"time"

	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/ctxsync"
	"github.com/siddhant2408/golang-libraries/ctxutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
	"github.com/streadway/amqp"
)

// Producer represents an AMQP producer.
type Producer func(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error

// SimpleProducer is a simple producer.
//
// It support confirmation optionally.
type SimpleProducer struct {
	Channel ChannelGetter
	Confirm bool

	mu    ctxsync.Mutex
	chn   *amqp.Channel
	cfmCh <-chan amqp.Confirmation
}

// Produce implements Producer.
//
// If confirmation is enabled, it returns when the confirmation is received.
func (p *SimpleProducer) Produce(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) (err error) {
	span, spanFinish := startTraceChildSpan(&ctx, "simple_producer", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, tracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.SpanTypeMessageProducer)
	opentracing_ext.SpanKindProducer.Set(span)
	err = tracingutils.TraceSyncLockerCtx(ctx, &p.mu)
	if err != nil {
		return errors.Wrap(err, "lock")
	}
	defer p.mu.Unlock()
	if exchange != "" {
		setTraceSpanTag(span, "exchange", exchange)
	}
	if key != "" {
		setTraceSpanTag(span, "routing_key", key)
	}
	if len(msg.Headers) > 0 {
		setTraceSpanTag(span, "headers", fmt.Sprint(msg.Headers))
	}
	setTraceSpanTagBody(span, msg.Body)
	err = p.produce(ctx, exchange, key, mandatory, immediate, msg)
	if err != nil {
		_ = p.close()
		return wrapErrorProducer(err, exchange, key, msg)
	}
	return nil
}

func (p *SimpleProducer) produce(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	chn, err := p.getChannel(ctx)
	if err != nil {
		return errors.Wrap(err, "get channel")
	}
	err = p.publish(ctx, chn, exchange, key, mandatory, immediate, msg)
	if err != nil {
		return errors.Wrap(err, "publish")
	}
	err = p.confirm(ctx)
	if err != nil {
		return errors.Wrap(err, "confirm")
	}
	return nil
}

func (p *SimpleProducer) publish(ctx context.Context, chn *amqp.Channel, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) (err error) {
	span, spanFinish := startTraceChildSpan(&ctx, "simple_producer.publish", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, tracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.SpanTypeMessageProducer)
	opentracing_ext.SpanKindProducer.Set(span)
	err = chn.Publish(exchange, key, mandatory, immediate, msg)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (p *SimpleProducer) confirm(ctx context.Context) (err error) {
	if p.cfmCh == nil {
		return nil
	}
	span, spanFinish := startTraceChildSpan(&ctx, "simple_producer.confirm", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, tracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.SpanTypeMessageProducer)
	opentracing_ext.SpanKindProducer.Set(span)
	select {
	case cfm, ok := <-p.cfmCh:
		if !ok {
			return errors.New("channel closed")
		}
		if !cfm.Ack {
			return errors.New("negative confirmation")
		}
		return nil
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "")
	}
}

func (p *SimpleProducer) getChannel(ctx context.Context) (*amqp.Channel, error) {
	if p.chn != nil {
		return p.chn, nil
	}
	chn, err := p.Channel(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "open channel")
	}
	cfmCh, err := p.initConfirm(chn)
	if err != nil {
		return nil, errors.Wrap(err, "confirm")
	}
	p.chn = chn
	p.cfmCh = cfmCh
	return p.chn, nil
}

func (p *SimpleProducer) initConfirm(chn *amqp.Channel) (<-chan amqp.Confirmation, error) {
	if !p.Confirm {
		return nil, nil
	}
	err := chn.Confirm(false)
	if err != nil {
		return nil, errors.Wrap(err, "set mode")
	}
	cfmCh := make(chan amqp.Confirmation, 1)
	chn.NotifyPublish(cfmCh)
	return cfmCh, nil
}

// Close closes the SimpleProducer.
//
// It closes the underlying channel.
func (p *SimpleProducer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.close()
}

func (p *SimpleProducer) close() error {
	chn := p.chn
	p.chn = nil
	p.cfmCh = nil
	if chn != nil {
		err := chn.Close()
		if err != nil {
			return errors.Wrap(err, "close channel")
		}
	}
	return nil
}

// TimeoutProducer sets a timeout to a Producer context.
type TimeoutProducer struct {
	Producer
	Timeout time.Duration
}

// Produce implements Producer.
func (p *TimeoutProducer) Produce(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) (err error) {
	ctx, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()
	return p.Producer(ctx, exchange, key, mandatory, immediate, msg)
}

// MultiProducer is a Producer that forwards to several Producers.
//
// Each Producer runs in a goroutine.
type MultiProducer struct {
	chn chan Producer
}

// NewMultiProducer returns a new MultiProducer.
//
// It starts the goroutines for the producers.
func NewMultiProducer(ps []Producer) *MultiProducer {
	chn := make(chan Producer, len(ps))
	for _, p := range ps {
		chn <- p
	}
	return &MultiProducer{
		chn: chn,
	}
}

// Produce implements Producer.
//
// The call is blocking until the message is produced.
func (mp *MultiProducer) Produce(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	select {
	case p := <-mp.chn:
		defer func() {
			mp.chn <- p
		}()
		return p(ctx, exchange, key, mandatory, immediate, msg)
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "")
	}
}

// BufferedProducer is a Producer with a buffer.
//
// It accumulates messages and sends them asynchronously.
type BufferedProducer struct {
	p       Producer
	ch      chan *bufferedProducerCall
	errFunc func(context.Context, error)
}

type bufferedProducerCall struct {
	exchange  string
	key       string
	mandatory bool
	immediate bool
	msg       amqp.Publishing
}

// NewBufferProducer returns a new BufferProducer.
func NewBufferProducer(p Producer, size int, errFunc func(context.Context, error)) *BufferedProducer {
	return &BufferedProducer{
		p:       p,
		ch:      make(chan *bufferedProducerCall, size),
		errFunc: errFunc,
	}
}

// Produce implements Producer.
//
// It writes the message to the buffer, which is processed asynchronously.
// If the buffer is full, it returns an error.
func (p *BufferedProducer) Produce(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	c := &bufferedProducerCall{
		exchange:  exchange,
		key:       key,
		mandatory: mandatory,
		immediate: immediate,
		msg:       msg,
	}
	select {
	case p.ch <- c:
		return nil
	default:
		err := errors.New("buffer full")
		err = wrapErrorValue(err, "buffer.size", cap(p.ch))
		err = wrapErrorProducer(err, exchange, key, msg)
		return err
	}
}

// Run runs the BufferedProducer.
//
// It reads the messages from the buffer and sends them to the sub Producer.
// If an error is returned by the sub Producer, the error function is called.
// The call is blocking until the Context is canceled.
func (p *BufferedProducer) Run(ctx context.Context) {
	for !ctxutils.IsDone(ctx) {
		select {
		case c := <-p.ch:
			p.produce(c)
		case <-ctx.Done():
			return
		}
	}
}

// Drain drains the BufferedProducer.
//
// It reads remaining messages from the buffer and sends them to the sub Producer.
// If an error is returned by the sub Producer, the error function is called.
// The call is blocking until the buffer is empty or the Context is canceled.
func (p *BufferedProducer) Drain(ctx context.Context) {
	for !ctxutils.IsDone(ctx) {
		select {
		case c := <-p.ch:
			p.produce(c)
		case <-ctx.Done():
			return
		default:
			return
		}
	}
}

func (p *BufferedProducer) produce(c *bufferedProducerCall) {
	ctx := context.Background()
	err := p.p(ctx, c.exchange, c.key, c.mandatory, c.immediate, c.msg)
	if err != nil {
		err = errors.Wrap(err, "AMQP buffered producer")
		p.errFunc(ctx, err)
	}
}

// ErrorProducer is a Producer that handles errors instead of returning them.
//
// It doesn't implement the Producer interface.
type ErrorProducer struct {
	p       Producer
	errFunc func(context.Context, error)
}

// NewErrorProducer returns a new ErrorProducer.
func NewErrorProducer(p Producer, errFunc func(context.Context, error)) *ErrorProducer {
	return &ErrorProducer{
		p:       p,
		errFunc: errFunc,
	}
}

// Produce produces the message to the sub Producer.
//
// If the sub Producer returns an error, the error handling function is called.
func (p *ErrorProducer) Produce(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) {
	err := p.p(ctx, exchange, key, mandatory, immediate, msg)
	if err != nil {
		err = errors.Wrap(err, "AMQP error producer")
		p.errFunc(ctx, err)
	}
}

// NewMultiConfirmProducer is a helper that creates several SimpleProducer with the confirm option, wrapped with a MultiProducer.
func NewMultiConfirmProducer(cg ChannelGetter, count int) (Producer, closeutils.WithOnErr) {
	ps := make([]Producer, count)
	cls := make([]closeutils.Err, count)
	for i := 0; i < count; i++ {
		p := &SimpleProducer{
			Channel: cg,
			Confirm: true,
		}
		ps[i] = p.Produce
		cls[i] = p.Close
	}
	mp := NewMultiProducer(ps)
	cl := func(oe closeutils.OnErr) {
		for i, cl := range cls {
			err := cl()
			if err != nil {
				err = errors.Wrapf(err, "simple: %d", i)
				oe(err)
			}
		}
	}
	return mp.Produce, cl
}

// Reproduce re-produces a delivered message.
// It removes the header "x-death".
func Reproduce(ctx context.Context, p Producer, exchange, key string, mandatory, immediate bool, dlv amqp.Delivery) error {
	pbl := deliveryToPublishing(dlv)
	return p(ctx, exchange, key, mandatory, immediate, pbl)
}

func wrapErrorProducer(err error, exchange string, key string, msg amqp.Publishing) error {
	if exchange != "" {
		err = wrapErrorValue(err, "exchange", exchange)
	}
	if key != "" {
		err = wrapErrorValue(err, "routing_key", key)
	}
	if len(msg.Headers) > 0 {
		err = wrapErrorValue(err, "headers", msg.Headers)
	}
	err = wrapErrorValueBody(err, msg.Body)
	return err
}
