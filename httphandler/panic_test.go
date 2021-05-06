package httphandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/siddhant2408/golang-libraries/panichandle"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestPanicNormal(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
	h := &Panic{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {}),
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	rec.Flush()
}

func TestPanicAbort(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/abort", nil)
	h := &Panic{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			panic(http.ErrAbortHandler)
		}),
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	rec.Flush()
}

func TestPanicPanic(t *testing.T) {
	var called testutils.CallCounter
	originalPH := panichandle.Handler
	defer func() {
		panichandle.Handler = originalPH
	}()
	panichandle.Handler = func(r interface{}) {
		called.Call()
		if r == nil {
			t.Fatal("no panic")
		}
	}
	h := &Panic{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			panic("lol")
		}),
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/panic", nil)
	h.ServeHTTP(rec, req)
	rec.Flush()
	called.AssertCalled(t)
}
