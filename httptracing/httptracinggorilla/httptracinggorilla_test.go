package httptracinggorilla

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestWrapRouter(t *testing.T) {
	var hCalled testutils.CallCounter
	h := func(w http.ResponseWriter, req *http.Request) {
		hCalled.Call()
		ctx := req.Context()
		span := opentracing.SpanFromContext(ctx)
		if span == nil {
			t.Fatal("no span")
		}
	}
	r := mux.NewRouter()
	r.Path("/test").HandlerFunc(h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost/test", nil)
	wr := WrapRouter(r)
	wr.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", w.Code, http.StatusOK)
	}
	hCalled.AssertCalled(t)
}

func TestWrapRouterNotFound(t *testing.T) {
	r := mux.NewRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost/test", nil)
	wr := WrapRouter(r)
	wr.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusNotFound {
		t.Fatalf("unexpected status code: got %d, want %d", w.Code, http.StatusNotFound)
	}
}
