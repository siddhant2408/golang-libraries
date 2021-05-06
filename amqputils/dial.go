package amqputils

import (
	"context"
	"math/rand"

	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/siddhant2408/golang-libraries/ctxutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
	"github.com/streadway/amqp"
)

// Dial dials with an URL.
func Dial(ctx context.Context, u string) (conn *amqp.Connection, err error) {
	span, spanFinish := startTraceChildSpan(&ctx, "dial", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, tracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.AppTypeRPC)
	opentracing_ext.SpanKindRPCClient.Set(span)
	setTraceSpanTag(span, "url", u)
	conn, err = amqp.Dial(u)
	err = errors.Wrap(err, "")
	return conn, err
}

// DialConfig dials with an URL and a config.
func DialConfig(ctx context.Context, u string, cfg amqp.Config) (conn *amqp.Connection, err error) {
	span, spanFinish := startTraceChildSpan(&ctx, "dial_config", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, tracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.AppTypeRPC)
	opentracing_ext.SpanKindRPCClient.Set(span)
	setTraceSpanTag(span, "url", u)
	conn, err = amqp.DialConfig(u, cfg)
	err = errors.Wrap(err, "")
	return conn, err
}

const dialNoConnectionAvailableErrMsg = "no connection available"

// URLsDialer dials a list of URLs.
type URLsDialer struct {
	URLs    []string
	DialURL func(context.Context, string) (*amqp.Connection, error)
}

// Dial dials the list of URLs.
// It attempts to dial each URL successively (in random order).
// If the connection is established, it returns it.
// Otherwise it goes to the next URL.
// If no URL succeeded, it returns the latest error.
func (d *URLsDialer) Dial(ctx context.Context) (conn *amqp.Connection, err error) {
	_, spanFinish := startTraceChildSpan(&ctx, "urls_dialer", &err)
	defer spanFinish()
	for _, u := range d.getURLs() {
		if ctxutils.IsDone(ctx) {
			return nil, errors.Wrap(ctx.Err(), "")
		}
		conn, err = d.DialURL(ctx, u)
		if err == nil {
			return conn, nil
		}
	}
	if err == nil {
		err = errors.New(dialNoConnectionAvailableErrMsg)
	} else {
		err = errors.Wrap(err, dialNoConnectionAvailableErrMsg)
	}
	err = wrapErrorValue(err, "urls", d.URLs)
	return nil, err
}

func (d *URLsDialer) getURLs() []string {
	us := make([]string, len(d.URLs))
	for i, j := range rand.Perm(len(d.URLs)) {
		us[i] = d.URLs[j]
	}
	return us
}
