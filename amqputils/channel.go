package amqputils

import (
	"context"

	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/siddhant2408/golang-libraries/ctxsync"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
	"github.com/streadway/amqp"
)

// ChannelGetter returns a new Channel.
type ChannelGetter func(context.Context) (*amqp.Channel, error)

// NewChannelGetterConnection returns a ChannelGetter for a Connection.
func NewChannelGetterConnection(conn *amqp.Connection) ChannelGetter {
	return func(ctx context.Context) (chn *amqp.Channel, err error) {
		return connectionChannel(ctx, conn)
	}
}

func connectionChannel(ctx context.Context, conn *amqp.Connection) (chn *amqp.Channel, err error) {
	span, spanFinish := startTraceChildSpan(&ctx, "connection.channel", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, tracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.AppTypeRPC)
	opentracing_ext.SpanKindRPCClient.Set(span)
	chn, err = conn.Channel()
	err = errors.Wrap(err, "")
	return chn, err
}

// ChannelPool is a pool of channel.
type ChannelPool struct {
	Channel ChannelGetter

	mu   ctxsync.Mutex
	chns []*amqp.Channel
}

// Get returns a channel from the pool.
// It opens a new one if the pool is empty.
// Warning: it is not guaranteed that the returned channel is open.
func (cp *ChannelPool) Get(ctx context.Context) (chn *amqp.Channel, err error) {
	_, spanFinish := startTraceChildSpan(&ctx, "channel_pool.get", &err)
	defer spanFinish()
	err = tracingutils.TraceSyncLockerCtx(ctx, &cp.mu)
	if err != nil {
		return nil, errors.Wrap(err, "lock")
	}
	defer cp.mu.Unlock()
	l := len(cp.chns)
	if l == 0 {
		chn, err = cp.Channel(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "open channel")
		}
		return chn, nil
	}
	i := l - 1
	chn = cp.chns[i]
	cp.chns[i] = nil
	cp.chns = cp.chns[:i]
	return chn, nil
}

// Put puts a channel to the pool.
func (cp *ChannelPool) Put(chn *amqp.Channel) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.chns = append(cp.chns, chn)
}

// Run runs a given function with a channel from the pool.
//
// When the given function returns, the channel is recycled to the pool, unless a *amqp.Error is returned.
// For other types of error, the channel is recycled to the pool.
func (cp *ChannelPool) Run(ctx context.Context, f func(context.Context, *amqp.Channel) error) (err error) {
	_, spanFinish := startTraceChildSpan(&ctx, "channel_pool.run", &err)
	defer spanFinish()
	ch, err := cp.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "get channel")
	}
	err = f(ctx, ch)
	if err != nil {
		errc := errors.UnwrapAll(err)
		if _, ok := errc.(*amqp.Error); ok {
			_ = ch.Close()
		} else {
			cp.Put(ch)
		}
		return err
	}
	cp.Put(ch)
	return nil
}

// Close closes the pool and all channels.
//
// It is OK to reuse the ChannelPool after this call.
func (cp *ChannelPool) Close() error {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	var firstErr error
	for _, chn := range cp.chns {
		err := chn.Close()
		if err != nil && firstErr == nil {
			firstErr = errors.Wrap(err, "close channel")
		}
	}
	cp.chns = nil
	return firstErr
}
