package httputils

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestReadRequestBody(t *testing.T) {
	ctx := context.Background()
	b1 := []byte("test")
	body := bytes.NewReader(b1)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", body)
	b2, err := ReadRequestBody(ctx, req)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if !bytes.Equal(b2, b1) {
		t.Fatal("not equal")
	}
}

func TestCopyRequestBody(t *testing.T) {
	ctx := context.Background()
	b1 := []byte("test")
	body := bytes.NewReader(b1)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", body)
	w := io.Discard
	written, err := CopyRequestBody(ctx, req, w)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if written != 4 {
		t.Fatalf("unexpected written: got %d, want %d", written, 4)
	}
}

func TestWriteResponse(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()
	b1 := []byte("test")
	WriteResponse(ctx, w, http.StatusOK, b1)
	w.Flush()
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusOK)
	}
	b2 := w.Body.Bytes()
	if !bytes.Equal(b2, b1) {
		t.Fatal("not equal")
	}
}

func TestWriteResponseText(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()
	s1 := "test"
	WriteResponseText(ctx, w, http.StatusOK, s1)
	w.Flush()
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusOK)
	}
	s2 := w.Body.String()
	if s2 != s1 {
		t.Fatal("not equal")
	}
}

func TestCopyResponse(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()
	b1 := []byte("test")
	r := bytes.NewReader(b1)
	CopyResponse(ctx, w, http.StatusOK, r)
	w.Flush()
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusOK)
	}
	b2 := w.Body.Bytes()
	if !bytes.Equal(b2, b1) {
		t.Fatal("not equal")
	}
}
