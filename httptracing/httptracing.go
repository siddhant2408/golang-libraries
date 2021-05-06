// Package httptracing provides HTTP tracing related utilities.
package httptracing

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/siddhant2408/golang-libraries/httpclientip"
	"github.com/siddhant2408/golang-libraries/httpurl"
	"github.com/siddhant2408/golang-libraries/tracingutils"
	ddtrace_ext "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
)

// Serve serves an HTTP request.
func Serve(h http.Handler, w http.ResponseWriter, req *http.Request, resource string) {
	ctx := req.Context()
	span, spanFinish := tracingutils.StartRootSpan(&ctx, "http.request", nil)
	defer spanFinish()
	tracingutils.SetSpanType(span, tracingutils.SpanTypeWeb)
	opentracing_ext.SpanKindRPCServer.Set(span)
	tracingutils.SetSpanResourceName(span, resource)
	span.SetTag(ddtrace_ext.HTTPMethod, req.Method)
	span.SetTag(ddtrace_ext.HTTPURL, httpurl.Get(req).String())
	span.SetTag("http.user_agent", req.UserAgent())
	clientIP, err := httpclientip.GetFromRequest(req)
	if err == nil {
		span.SetTag("http.client_ip", clientIP.String())
	}
	req = req.WithContext(ctx)
	srw := &statusResponseWriter{
		ResponseWriter: w,
		span:           span,
	}
	h.ServeHTTP(srw, req)
	srw.ensureStatus()
}

type statusResponseWriter struct {
	http.ResponseWriter
	span   opentracing.Span
	status int
}

func (w *statusResponseWriter) Write(b []byte) (int, error) {
	w.ensureStatus()
	return w.ResponseWriter.Write(b)
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
	w.span.SetTag(ddtrace_ext.HTTPCode, status)
	if status >= 500 && status < 600 {
		tracingutils.SetSpanHasError(w.span)
	}
}

func (w *statusResponseWriter) ensureStatus() {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
}

// ResourceResolver resolves a resource for a Request.
type ResourceResolver func(*http.Request) string

// WrapHandler wraps a handler.
func WrapHandler(h http.Handler, resourceResolver ResourceResolver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		resource := resourceResolver(req)
		Serve(h, w, req, resource)
	})
}

// WrapServeMux wraps a ServeMux.
func WrapServeMux(m *http.ServeMux) http.Handler {
	return WrapHandler(m, NewServeMuxResourceResolver(m))
}

// NewServeMuxResourceResolver returns a new ResourceResolver for a ServeMux.
func NewServeMuxResourceResolver(m *http.ServeMux) ResourceResolver {
	return func(req *http.Request) string {
		_, route := m.Handler(req)
		if route == "" {
			route = "unknown"
		}
		resource := req.Method + " " + route
		return resource
	}
}

const (
	roundTripperExternalServiceName = "go-http"
)

type roundTripper struct {
	http.RoundTripper
}

// WrapRoundTripper wraps a http.RoundTripper.
func WrapRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &roundTripper{
		RoundTripper: rt,
	}
}

func (rt *roundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	ctx := req.Context()
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "http.roundtrip", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, roundTripperExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.SpanTypeHTTP)
	opentracing_ext.SpanKindRPCClient.Set(span)
	span.SetTag(ddtrace_ext.HTTPMethod, req.Method)
	span.SetTag(ddtrace_ext.HTTPURL, req.URL.String())
	req = req.WithContext(ctx)
	resp, err = rt.RoundTripper.RoundTrip(req)
	if err != nil {
		// DO NOT wrap the error here.
		// It causes bug if the parent error expect a specific error type (implementing an interface).
		return nil, err
	}
	span.SetTag(ddtrace_ext.HTTPCode, resp.StatusCode)
	if resp.StatusCode >= 500 && resp.StatusCode < 600 {
		tracingutils.SetSpanHasError(span)
	}
	return resp, nil
}

// WrapDefaultTransport initializes http.DefaultTransport.
// It wraps it with WrapRoundTripper.
func WrapDefaultTransport() {
	http.DefaultTransport = WrapRoundTripper(http.DefaultTransport)
}
