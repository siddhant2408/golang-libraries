package amqputils

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

const (
	tracingExternalServiceName = "go-amqp"
)

func startTraceRootSpan(pctx *context.Context, op string, perr *error) (opentracing.Span, closeutils.F) {
	return tracingutils.StartRootSpan(pctx, "amqp."+op, perr)
}

func startTraceChildSpan(pctx *context.Context, op string, perr *error) (opentracing.Span, closeutils.F) {
	return tracingutils.StartChildSpan(pctx, "amqp."+op, perr)
}

func setTraceSpanTag(span opentracing.Span, key string, val interface{}) {
	span.SetTag("amqp."+key, val)
}

func setTraceSpanTagBody(span opentracing.Span, body []byte) {
	setTraceSpanTag(span, "body", bodyTruncateConvert(body))
}
