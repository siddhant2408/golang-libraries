package mongo

import (
	"context"
	"fmt"
	"reflect" //nolint:depguard // Used for JSON encoding.

	"github.com/opentracing/opentracing-go"
	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// TracingExternalServiceName is the tracing external service name.
	TracingExternalServiceName = "go-mongodb"
)

func startTraceSpan(pctx *context.Context, op string, perr *error, werr func(error) error) (opentracing.Span, closeutils.F) {
	span, spanFinish := tracingutils.StartChildSpan(pctx, "mongo."+op, perr)
	tracingutils.SetSpanServiceName(span, TracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.SpanTypeMongoDB)
	opentracing_ext.SpanKindRPCClient.Set(span)
	spanFinish = wrapTraceSpanFinishError(spanFinish, op, perr, werr)
	return span, spanFinish
}

func wrapTraceSpanFinishError(f func(), op string, perr *error, werr func(error) error) func() {
	if perr == nil {
		return f
	}
	return func() {
		wrapErrorReturn(op, perr, werr)
		f()
	}
}

func setTraceSpanTag(span opentracing.Span, key string, val interface{}) {
	span.SetTag("mongo."+key, val)
}

const (
	tracingSpanJSONMaxSize = 4 << 10 // 4 KiB
)

func setTraceSpanTagJSON(span opentracing.Span, key string, val interface{}) {
	// We don't do anything if the span is noop, because JSON marshal can be very costly.
	if tracingutils.IsSpanNoop(span) {
		return
	}
	b, err := marshalJSON(val)
	if err != nil {
		b = []byte(fmt.Sprintf("error: %v", err))
	}
	if len(b) > tracingSpanJSONMaxSize {
		b = b[:tracingSpanJSONMaxSize]
	}
	setTraceSpanTag(span, key, string(b))
}

func marshalJSON(v interface{}) ([]byte, error) {
	// bson.MarshalExtJSON can't marshal some types that are at the top level.
	// So we wrap them in an object, and remove this wrapping after the marshaling is done.
	topLevelHack := shouldUseMarshalJSONTopLevelHack(v)
	if topLevelHack {
		v = marshalJSONWrapper{
			A: v,
		}
	}
	b, err := bson.MarshalExtJSON(v, false, false)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if topLevelHack {
		start := 5        // `{"a":`
		end := len(b) - 1 // `}`
		b = b[start:end]
	}
	return b, nil
}

func shouldUseMarshalJSONTopLevelHack(v interface{}) bool {
	switch v.(type) {
	case bson.M, bson.D:
		return false
	case nil:
		return true
	}
	kd := unwrapPointerKind(v)
	switch kd {
	case reflect.Map, reflect.Struct:
		return false
	}
	return true
}

func unwrapPointerKind(v interface{}) reflect.Kind {
	typ := reflect.TypeOf(v)
	for {
		kd := typ.Kind()
		if kd != reflect.Ptr {
			return kd
		}
		typ = typ.Elem()
	}
}

type marshalJSONWrapper struct {
	A interface{} `json:"a"`
}
