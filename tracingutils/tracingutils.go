// Package tracingutils provides tracing related utilities.
package tracingutils

import (
	"context"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"
	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/errors"
	ddtrace_ext "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
)

// Span types.
const (
	AppTypeWeb              = ddtrace_ext.AppTypeWeb
	AppTypeDB               = ddtrace_ext.AppTypeDB
	AppTypeCache            = ddtrace_ext.AppTypeCache
	AppTypeRPC              = ddtrace_ext.AppTypeRPC
	SpanTypeWeb             = ddtrace_ext.SpanTypeWeb
	SpanTypeHTTP            = ddtrace_ext.SpanTypeHTTP
	SpanTypeSQL             = ddtrace_ext.SpanTypeSQL
	SpanTypeCassandra       = ddtrace_ext.SpanTypeCassandra
	SpanTypeRedis           = ddtrace_ext.SpanTypeRedis
	SpanTypeMemcached       = ddtrace_ext.SpanTypeMemcached
	SpanTypeMongoDB         = ddtrace_ext.SpanTypeMongoDB
	SpanTypeElasticSearch   = ddtrace_ext.SpanTypeElasticSearch
	SpanTypeLevelDB         = ddtrace_ext.SpanTypeLevelDB
	SpanTypeDNS             = ddtrace_ext.SpanTypeDNS
	SpanTypeMessageConsumer = ddtrace_ext.SpanTypeMessageConsumer
	SpanTypeMessageProducer = ddtrace_ext.SpanTypeMessageProducer
	SpanTypeConsul          = ddtrace_ext.SpanTypeConsul
)

// SetSpanType sets the span's type.
func SetSpanType(span opentracing.Span, typ string) {
	span.SetTag(ddtrace_ext.SpanType, typ)
}

// SetSpanServiceName sets the span's service name.
func SetSpanServiceName(span opentracing.Span, srv string) {
	span.SetTag(ddtrace_ext.ServiceName, srv)
}

// SetSpanResourceName sets the span's resource name.
func SetSpanResourceName(span opentracing.Span, res string) {
	span.SetTag(ddtrace_ext.ResourceName, res)
}

// SetSpanHasError sets the "has error" information on a span.
func SetSpanHasError(span opentracing.Span) {
	opentracing_ext.Error.Set(span, true)
}

// SetSpanError set the span's error.
func SetSpanError(span opentracing.Span, err error) {
	// We don't do anything if the span is noop, because formatting the error can be very costly.
	if IsSpanNoop(span) {
		return
	}
	SetSpanHasError(span)
	span.SetTag(ddtrace_ext.ErrorMsg, err.Error())
	span.SetTag(ddtrace_ext.ErrorStack, fmt.Sprintf("%+v", err))
}

// IsSpanNoop returns true if the span is noop, false otherwise.
func IsSpanNoop(span opentracing.Span) bool {
	_, ok := span.Tracer().(opentracing.NoopTracer)
	return ok
}

// StartRootSpan calls StartRootSpanWithTracer with the global tracer.
func StartRootSpan(pctx *context.Context, operationName string, perr *error) (opentracing.Span, closeutils.F) {
	return StartRootSpanWithTracer(pctx, opentracing.GlobalTracer(), operationName, perr)
}

// StartRootSpanWithTracer starts a root span.
// This span has no parent, and starts a new trace.
// It returns a "close" function that allows to finish the span.
// The pctx parameter updates the context with  new one containing the span.
// The perr parameter automatically tracks the returned error.
func StartRootSpanWithTracer(pctx *context.Context, tr opentracing.Tracer, operationName string, perr *error) (opentracing.Span, closeutils.F) {
	return startSpan(pctx, tr, operationName, nil, perr)
}

// StartChildSpan start a child span.
// If the current context has no span, a noop span is returned.
// It uses the same tracer as the current span.
// It returns a "close" function that allows to finish the span.
// The pctx parameter updates the context with  new one containing the span.
// The perr parameter automatically tracks the returned error.
func StartChildSpan(pctx *context.Context, operationName string, perr *error) (opentracing.Span, closeutils.F) {
	ctx := *pctx
	var tr opentracing.Tracer
	var opts []opentracing.StartSpanOption
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		tr = span.Tracer()
		opts = []opentracing.StartSpanOption{opentracing.ChildOf(span.Context())}
	} else {
		tr = opentracing.NoopTracer{}
	}
	return startSpan(pctx, tr, operationName, opts, perr)
}

func startSpan(pctx *context.Context, tr opentracing.Tracer, operationName string, opts []opentracing.StartSpanOption, perr *error) (opentracing.Span, closeutils.F) {
	span := tr.StartSpan(operationName, opts...)
	ctx := *pctx
	ctx = opentracing.ContextWithSpan(ctx, span)
	*pctx = ctx
	finish := newSpanFinish(span, perr)
	return span, finish
}

func newSpanFinish(span opentracing.Span, perr *error) func() {
	if perr == nil {
		return span.Finish
	}
	return func() {
		err := *perr
		if err != nil && !errors.IsIgnored(err) {
			SetSpanError(span, err)
		}
		span.Finish()
	}
}

// TraceSyncLocker traces sync.Locker.Lock.
func TraceSyncLocker(ctx context.Context, l interface{ Lock() }) {
	_, spanFinish := StartChildSpan(&ctx, "sync_locker", nil)
	defer spanFinish()
	l.Lock()
}

// TraceSyncLockerCtx traces a sync Locker with a context.
func TraceSyncLockerCtx(ctx context.Context, l interface{ LockCtx(context.Context) error }) error {
	_, spanFinish := StartChildSpan(&ctx, "sync_locker_ctx", nil)
	defer spanFinish()
	return l.LockCtx(ctx)
}
