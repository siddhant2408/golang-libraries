package errorhandle

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/httperrors"
)

func TestHandle(t *testing.T) {
	ctx := context.Background()
	err := errors.New("error")
	Handle(ctx, err)
}

func TestHandleDefault(t *testing.T) {
	ctx := context.Background()
	err := errors.New("error")
	HandleDefault(ctx, err)
}

func TestNewHandler(t *testing.T) {
	ctx := context.Background()
	err := errors.New("error")
	h := NewHandler()
	h(ctx, err)
}

func TestIgnored(t *testing.T) {
	ctx := context.Background()
	err := errors.New("error")
	err = errors.Ignore(err)
	Handle(ctx, err)
}

func TestWait(t *testing.T) {
	ctx := context.Background()
	err := errors.New("error")
	Handle(ctx, err, Wait())
}

func TestSentryID(t *testing.T) {
	ctx := context.Background()
	var sentryID string
	err := errors.New("error")
	Handle(ctx, err, SentryID(&sentryID))
	if sentryID == "" {
		t.Fatal("empty ID")
	}
}

func TestHTTPServerRequest(t *testing.T) {
	ctx := context.Background()
	err := errors.New("error")
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	err = httperrors.WithServerRequest(err, req)
	Handle(ctx, err)
}

func TestFilterRavenErrorWebCrawlers(t *testing.T) {
	err := errors.New(sentryErrorDropWebCrawlers)
	err = errors.Wrap(err, "test")
	filter := filterRavenError(err)
	if filter {
		t.Fatal("should be false")
	}
}

func TestFilterRavenErrorOther(t *testing.T) {
	err := errors.New("test")
	filter := filterRavenError(err)
	if !filter {
		t.Fatal("should be true")
	}
}

func TestTraceSpanTags(t *testing.T) {
	ctx := context.Background()
	tr := mocktracer.New()
	span := tr.StartSpan("test").(*mocktracer.MockSpan) //nolint:errcheck
	ctx = opentracing.ContextWithSpan(ctx, span)
	err := errors.New("error")
	Handle(ctx, err)
	sentryID, ok := span.Tag(traceSpanTagSentry).(string)
	if !ok || sentryID == "" {
		t.Fatal("missing trace span tag")
	}
}

func TestHTTPHeader(t *testing.T) {
	ctx := context.Background()
	err := errors.New("error")
	hr := make(http.Header)
	Handle(ctx, err, HTTPHeader(hr))
	sentryID := hr.Get(httpHeaderSentry)
	if sentryID == "" {
		t.Fatal("missing HTTP header")
	}
}
