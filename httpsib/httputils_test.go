package httpsib

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrapHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
	h = WrapHandler(
		h,
		MaxConcurrent(10, nil),
		RequestBodyMaxBytes(1<<10),
		ContextCancel(false),
	)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	rec.Flush()
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", rec.Code, http.StatusOK)
	}
}
