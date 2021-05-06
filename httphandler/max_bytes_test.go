package httphandler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestRequestBodyMaxBytes(t *testing.T) {
	body := "test"
	req := httptest.NewRequest(http.MethodPost, "http://example.com", strings.NewReader(body))
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			testutils.ErrorErr(t, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if string(b) != body {
			err := errors.Newf("unexpected body: got %q, want %q", string(b), body)
			testutils.ErrorErr(t, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, _ = io.WriteString(w, "OK")
	})
	h = &RequestBodyMaxBytes{
		Handler:  h,
		MaxBytes: 1 << 20,
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	rec.Flush()
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestRequestBodyMaxBytesErrorMaxBytes(t *testing.T) {
	body := "test"
	req := httptest.NewRequest(http.MethodPost, "http://example.com", strings.NewReader(body))
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, err := io.Copy(io.Discard, req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t.Error("no error")
	})
	h = &RequestBodyMaxBytes{
		Handler:  h,
		MaxBytes: 1,
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	rec.Flush()
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status code: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
