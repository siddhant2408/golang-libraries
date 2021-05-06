package tracingutils

import (
	"context"
	"sync"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/siddhant2408/golang-libraries/ctxsync"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	ddtrace_ext "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
)

func TestSetSpanType(t *testing.T) {
	span := startMockSpan()
	SetSpanType(span, SpanTypeWeb)
	if span.Tag(ddtrace_ext.SpanType) != SpanTypeWeb {
		t.Fatalf("unexpected tag %q: got %q, want %q", ddtrace_ext.SpanType, span.Tag(ddtrace_ext.SpanType), SpanTypeWeb)
	}
}

func TestSetSpanServiceName(t *testing.T) {
	span := startMockSpan()
	SetSpanServiceName(span, "name")
	if span.Tag(ddtrace_ext.ServiceName) != "name" {
		t.Fatalf("unexpected tag %q: got %q, want %q", ddtrace_ext.ServiceName, span.Tag(ddtrace_ext.SpanType), "name")
	}
}

func TestSetSpanResourceName(t *testing.T) {
	span := startMockSpan()
	SetSpanResourceName(span, "name")
	if span.Tag(ddtrace_ext.ResourceName) != "name" {
		t.Fatalf("unexpected tag %q: got %q, want %q", ddtrace_ext.ResourceName, span.Tag(ddtrace_ext.ResourceName), "name")
	}
}

func TestSetSpanError(t *testing.T) {
	span := startMockSpan()
	err := errors.New("error")
	SetSpanError(span, err)
	if span.Tag(string(opentracing_ext.Error)) != true {
		t.Fatalf("unexpected tag %q: got %v, want %t", opentracing_ext.Error, span.Tag(string(opentracing_ext.Error)), true)
	}
}

func TestIsSpanNoopTrue(t *testing.T) {
	span := startMockSpan()
	noop := IsSpanNoop(span)
	if noop {
		t.Fatalf("unexpected noop: got %t, want %t", noop, false)
	}
}

func TestIsSpanNoopFalse(t *testing.T) {
	span := opentracing.NoopTracer{}.StartSpan("test")
	noop := IsSpanNoop(span)
	if !noop {
		t.Fatalf("unexpected noop: got %t, want %t", noop, true)
	}
}

func TestStartRootSpan(t *testing.T) {
	ctx := context.Background()
	span, spanFinish := StartRootSpan(&ctx, "test", nil)
	defer spanFinish()
	if span == nil {
		t.Fatal("nil span")
	}
	if opentracing.SpanFromContext(ctx) != span {
		t.Fatal("context span not equal")
	}
}

func TestStartRootSpanError(t *testing.T) {
	ctx := context.Background()
	var err error
	span, spanFinish := StartRootSpan(&ctx, "test", &err)
	defer spanFinish()
	if span == nil {
		t.Fatal("nil span")
	}
	if opentracing.SpanFromContext(ctx) != span {
		t.Fatal("context span not equal")
	}
	err = errors.New("error")
}

func TestStartChildSpan(t *testing.T) {
	tr := mocktracer.New()
	span := tr.StartSpan("test")
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	span, spanFinish := StartChildSpan(&ctx, "test", nil)
	defer spanFinish()
	if IsSpanNoop(span) {
		t.Fatal("noop span")
	}
	if opentracing.SpanFromContext(ctx) != span {
		t.Fatal("context span not equal")
	}
}

func TestStartChildSpanNoop(t *testing.T) {
	ctx := context.Background()
	span, spanFinish := StartChildSpan(&ctx, "test", nil)
	defer spanFinish()
	if !IsSpanNoop(span) {
		t.Fatal("not noop span")
	}
	if opentracing.SpanFromContext(ctx) != span {
		t.Fatal("context span not equal")
	}
}

func TestTraceSyncLocker(t *testing.T) {
	ctx := context.Background()
	var mu sync.Mutex
	TraceSyncLocker(ctx, &mu)
}

func TestTraceSyncLockerCtx(t *testing.T) {
	ctx := context.Background()
	var mu ctxsync.Mutex
	err := TraceSyncLockerCtx(ctx, &mu)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func startMockSpan() *mocktracer.MockSpan {
	tr := mocktracer.New()
	span := tr.StartSpan("test").(*mocktracer.MockSpan) //nolint:errcheck
	return span
}
