package httptracing

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/httpclientip"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestServe(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		span := opentracing.SpanFromContext(ctx)
		if span == nil {
			t.Fatal("no span")
		}
		_, _ = w.Write([]byte("test"))
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req = req.WithContext(httpclientip.SetToContext(req.Context(), net.ParseIP("127.0.0.1")))
	rsc := "test"
	Serve(h, w, req, rsc)
	w.Flush()
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestServe500(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req = req.WithContext(httpclientip.SetToContext(req.Context(), net.ParseIP("127.0.0.1")))
	rsc := "test"
	Serve(h, w, req, rsc)
	w.Flush()
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestWrapServeMux(t *testing.T) {
	var hCalled testutils.CallCounter
	h := func(w http.ResponseWriter, req *http.Request) {
		hCalled.Call()
		ctx := req.Context()
		span := opentracing.SpanFromContext(ctx)
		if span == nil {
			t.Fatal("no span")
		}
	}
	m := http.NewServeMux()
	m.HandleFunc("/test", h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost/test", nil)
	wm := WrapServeMux(m)
	wm.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", w.Code, http.StatusOK)
	}
	hCalled.AssertCalled(t)
}

func TestWrapServeMuxNotFound(t *testing.T) {
	m := http.NewServeMux()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost/test", nil)
	wm := WrapServeMux(m)
	wm.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusNotFound {
		t.Fatalf("unexpected status code: got %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestWrapRoundTripper(t *testing.T) {
	rt := &testRoundTripper{
		resp: &http.Response{},
	}
	wrt := WrapRoundTripper(rt)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	resp, err := wrt.RoundTrip(req) //nolint:bodyclose // It's a fake response.
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if resp == nil {
		t.Fatal("no response")
	}
}

func TestWrapRoundTripper500(t *testing.T) {
	rt := &testRoundTripper{
		resp: &http.Response{
			StatusCode: http.StatusInternalServerError,
		},
	}
	wrt := WrapRoundTripper(rt)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	resp, err := wrt.RoundTrip(req) //nolint:bodyclose // It's a fake response.
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if resp == nil {
		t.Fatal("no response")
	}
}

func TestWrapRoundTripperError(t *testing.T) {
	rt := &testRoundTripper{
		err: errors.New("error"),
	}
	wrt := WrapRoundTripper(rt)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	_, err := wrt.RoundTrip(req) //nolint:bodyclose // No response is returned.
	if err == nil {
		t.Fatal("no error")
	}
}

type testRoundTripper struct {
	resp *http.Response
	err  error
}

func (rt *testRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	return rt.resp, rt.err
}
