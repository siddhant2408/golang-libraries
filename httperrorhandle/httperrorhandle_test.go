package httperrorhandle

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestHandle(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	err := errors.New("error")
	res := Handle(ctx, w, req, err)
	if res.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected code: got %d, want %d", res.Code, http.StatusInternalServerError)
	}
	if res.Text == "" {
		t.Fatal("empty")
	}
}

func TestServe(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	var fCalled testutils.CallCounter
	f := func(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
		fCalled.Call()
		return nil
	}
	Serve(w, req, f, nil)
	fCalled.AssertCalled(t)
}

func TestServeError(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	f := func(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
		return errors.New("error")
	}
	var onErrCalled testutils.CallCounter
	onErr := func(ctx context.Context, w http.ResponseWriter, req *http.Request, err error, hres *HandleResult) {
		onErrCalled.Call()
	}
	Serve(w, req, f, onErr)
	onErrCalled.AssertCalled(t)
}
