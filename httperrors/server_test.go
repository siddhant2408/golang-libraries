package httperrors

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/httpclientip"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestServerRequest(t *testing.T) {
	err := errors.New("error")
	req1 := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	err = WithServerRequest(err, req1)
	err = errors.Wrap(err, "test")
	req2 := GetServerRequest(err)
	if req2 != req1 {
		t.Fatal("not equal")
	}
}

func TestWithServerRequestNil(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	err := WithServerRequest(nil, req)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestGetServerRequestNil(t *testing.T) {
	err := errors.New("error")
	req := GetServerRequest(err)
	if req != nil {
		t.Fatal("not nil")
	}
}

func TestServerRequestError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	err := errors.New("error")
	err = WithServerRequest(err, req)
	msg := err.Error()
	expectedMsg := "HTTP server request GET http://localhost: error"
	if msg != expectedMsg {
		t.Fatalf("unexpected message: got %q, want %q", msg, expectedMsg)
	}
}

func BenchmarkServerRequestError(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	err := errors.New("error")
	err = WithServerRequest(err, req)
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(io.Discard, "%v", err)
	}
}

func TestServerRequestFormat(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	h := make(http.Header)
	h.Set("User-Agent", "test")
	req.Header = h
	ip := net.ParseIP("1.2.3.4")
	req = req.WithContext(httpclientip.SetToContext(req.Context(), ip))
	err := errors.New("error")
	err = WithServerRequest(err, req)
	msg := fmt.Sprintf("%+v", err)
	expectedRegexp := regexp.MustCompile(`^HTTP server request
	method: GET
	URL: http:\/\/localhost
	client IP: 1\.2\.3\.4
	user-agent: test
	headers:
		User-Agent: test\n
stack
(\t.+\n)+error$`)
	if !expectedRegexp.MatchString(msg) {
		t.Fatalf("unexpected message:\ngot: %q\nwant match: %q", msg, expectedRegexp)
	}
}

func BenchmarkServerRequestFormat(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	h := make(http.Header)
	h.Set("User-Agent", "test")
	req.Header = h
	ip := net.ParseIP("1.2.3.4")
	req = req.WithContext(httpclientip.SetToContext(req.Context(), ip))
	err := errors.New("error")
	err = WithServerRequest(err, req)
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(io.Discard, "%+v", err)
	}
}

func TestServerCode(t *testing.T) {
	err := errors.New("error")
	err = WithServerCode(err, http.StatusBadRequest)
	code, text := GetServerCodeText(err)
	if code != http.StatusBadRequest {
		t.Fatalf("unexpected code: got %d, want %d", code, http.StatusBadRequest)
	}
	expectedText := "error"
	if text != expectedText {
		t.Fatalf("unexpected text: got %q, want %q", text, expectedText)
	}
}

func TestWithServerCodeNil(t *testing.T) {
	err := WithServerCode(nil, http.StatusBadRequest)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestGetServerCodeTextInternal(t *testing.T) {
	err := errors.New("error")
	code, text := GetServerCodeText(err)
	if code != http.StatusInternalServerError {
		t.Fatalf("unexpected code: got %d, want %d", code, http.StatusInternalServerError)
	}
	expectedText := "Internal Server Error"
	if text != expectedText {
		t.Fatalf("unexpected text: got %q, want %q", text, expectedText)
	}
}

func TestServerCodeFormat(t *testing.T) {
	err := errors.New("error")
	err = WithServerCode(err, http.StatusBadRequest)
	msg := fmt.Sprint(err)
	expectedMsg := "HTTP server code 400: error"
	if msg != expectedMsg {
		t.Fatalf("unexpected message: got %q, want %q", msg, expectedMsg)
	}
}

func BenchmarkServerCodeFormat(b *testing.B) {
	err := errors.New("error")
	err = WithServerCode(err, http.StatusBadRequest)
	for i := 0; i < b.N; i++ {
		fmt.Fprint(io.Discard, err)
	}
}
