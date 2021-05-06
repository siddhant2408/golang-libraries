package httphandler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/_ping", nil)
	h := &Ping{}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	rec.Flush()
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestPingFallback(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
	called := false
	h := &Ping{
		Handler: http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			called = true
		}),
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	rec.Flush()
	if !called {
		t.Fatal("not called")
	}
}
