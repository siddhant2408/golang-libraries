package httphandler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/siddhant2408/golang-libraries/ctxutils"
)

func TestNoContextCancel(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	ctx := req.Context()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	req = req.WithContext(ctx)
	h := &NoContextCancel{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			done := ctxutils.IsDone(req.Context())
			if done {
				t.Fatal("done")
			}
		}),
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	rec.Flush()
}
