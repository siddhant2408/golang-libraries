package httphandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestMaxConcurrent(t *testing.T) {
	var thCalled testutils.CallCounter
	th := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		thCalled.Call()
	})
	h := &MaxConcurrent{
		Handler: th,
		Limit:   10,
	}
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusOK)
	}
	thCalled.AssertCalled(t)
}

func TestMaxConcurrentLimitReached(t *testing.T) {
	h := &MaxConcurrent{
		Limit:   10,
		counter: 10, // Yes it's bad, but testing it for real is too complex.
	}
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusServiceUnavailable)
	}
}

func TestMaxConcurrentLimitReachedHandler(t *testing.T) {
	var lrhCalled testutils.CallCounter
	lrh := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		lrhCalled.Call()
	})
	h := &MaxConcurrent{
		Limit:               10,
		LimitReachedHandler: lrh,
		counter:             10, // Yes it's bad, but testing it for real is too complex.
	}
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	w.Flush()
	lrhCalled.AssertCalled(t)
}
