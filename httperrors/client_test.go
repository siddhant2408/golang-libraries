package httperrors

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestClientRequest(t *testing.T) {
	req1 := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	err := errors.New("error")
	err = WithClientRequest(err, req1)
	req2 := GetClientRequest(err)
	if req2 != req1 {
		t.Fatal("not equal")
	}
}

func TestWithClientRequestNil(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	err := WithClientRequest(nil, req)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestGetClientRequestNil(t *testing.T) {
	err := errors.New("error")
	req := GetClientRequest(err)
	if req != nil {
		t.Fatal("not nil")
	}
}

func TestClientRequestError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	err := errors.New("error")
	err = WithClientRequest(err, req)
	msg := err.Error()
	expectedMsg := "HTTP client request GET http://localhost: error"
	if msg != expectedMsg {
		t.Fatalf("unexpected message: got %q, want %q", msg, expectedMsg)
	}
}

func BenchmarkClientRequestError(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	err := errors.New("error")
	err = WithClientRequest(err, req)
	for i := 0; i < b.N; i++ {
		fmt.Fprint(io.Discard, err)
	}
}

func TestClientRequestFormat(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	err := errors.New("error")
	err = WithClientRequest(err, req)
	msg := fmt.Sprintf("%+v", err)
	expectedRegexp := regexp.MustCompile(`^HTTP client request
	method: GET
	URL: http:\/\/localhost\n
stack
(\t.+\n)+error`)
	if !expectedRegexp.MatchString(msg) {
		t.Fatalf("unexpected message:\ngot: %q\nwant match: %q", msg, expectedRegexp)
	}
}

func BenchmarkClientRequestFormat(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	err := errors.New("error")
	err = WithClientRequest(err, req)
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(io.Discard, "%+v", err)
	}
}

func TestClientResponse(t *testing.T) {
	cr1 := &ClientResponse{
		Response: &http.Response{},
		Body:     []byte("test"),
	}
	err := errors.New("error")
	err = WithClientResponse(err, cr1)
	cr2 := GetClientResponse(err)
	testutils.Compare(t, "unexpected client response", cr2, cr1)
}

func TestClientResponseNil(t *testing.T) {
	cr := &ClientResponse{
		Response: &http.Response{},
		Body:     []byte("test"),
	}
	err := WithClientResponse(nil, cr)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestGetClientResponseNil(t *testing.T) {
	err := errors.New("error")
	cr := GetClientResponse(err)
	if cr != nil {
		t.Fatal("not nil")
	}
}

func TestClientResponseError(t *testing.T) {
	cr := &ClientResponse{
		Response: &http.Response{
			Status:     "400 Bad Request",
			StatusCode: http.StatusBadRequest,
		},
		Body: []byte("test"),
	}
	err := errors.New("error")
	err = WithClientResponse(err, cr)
	msg := err.Error()
	expectedMsg := "HTTP client response 400 Bad Request: error"
	if msg != expectedMsg {
		t.Fatalf("unexpected message: got %q, want %q", msg, expectedMsg)
	}
}

func BenchmarkClientResponseError(b *testing.B) {
	cr := &ClientResponse{
		Response: &http.Response{
			Status:     "400 Bad Request",
			StatusCode: http.StatusBadRequest,
		},
		Body: []byte("test"),
	}
	err := errors.New("error")
	err = WithClientResponse(err, cr)
	for i := 0; i < b.N; i++ {
		fmt.Fprint(io.Discard, err)
	}
}

func TestClientResponseFormatBodyString(t *testing.T) {
	body := []byte("test")
	cr := &ClientResponse{
		Response: &http.Response{
			Status:     "400 Bad Request",
			StatusCode: http.StatusBadRequest,
			Proto:      "HTTP/1.1",
			Header: http.Header{
				"Foo": {"bar"},
			},
			ContentLength: int64(len(body)),
		},
		Body: body,
	}
	err := errors.New("error")
	err = WithClientResponse(err, cr)
	msg := fmt.Sprintf("%+v", err)
	expectedRegexp := regexp.MustCompile(`^HTTP client response
	status: 400 Bad Request
	proto: HTTP/1\.1
	content length: 4
	headers:
		Foo: bar
	body:
================ begin ================
test
================= end =================

stack
(\t.+\n)+error$`)
	if !expectedRegexp.MatchString(msg) {
		t.Fatalf("unexpected message:\ngot: %q\nwant match: %q", msg, expectedRegexp)
	}
}

func BenchmarkClientResponseFormatBodyString(b *testing.B) {
	body := []byte("test")
	cr := &ClientResponse{
		Response: &http.Response{
			Status:     "400 Bad Request",
			StatusCode: http.StatusBadRequest,
			Proto:      "HTTP/1.1",
			Header: http.Header{
				"Foo": {"bar"},
			},
			ContentLength: int64(len(body)),
		},
		Body: body,
	}
	err := errors.New("error")
	err = WithClientResponse(err, cr)
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(io.Discard, "%+v", err)
	}
}

func TestClientResponseFormatBodyStringTruncated(t *testing.T) {
	body := bytes.Repeat([]byte("test"), clientResponseBodyMaxSizeString)
	cr := &ClientResponse{
		Response: &http.Response{
			Status:     "400 Bad Request",
			StatusCode: http.StatusBadRequest,
			Proto:      "HTTP/1.1",
			Header: http.Header{
				"Foo": {"bar"},
			},
			ContentLength: int64(len(body)),
		},
		Body: body,
	}
	err := errors.New("error")
	err = WithClientResponse(err, cr)
	msg := fmt.Sprintf("%+v", err)
	expectedRegexp := regexp.MustCompile(`^HTTP client response
	status: 400 Bad Request
	proto: HTTP/1\.1
	content length: 16384
	headers:
		Foo: bar
	body:
================ begin ================
(test)+
\(truncated to 4096 bytes\)
================= end =================

stack
(\t.+\n)+error$`)
	if !expectedRegexp.MatchString(msg) {
		t.Fatalf("unexpected message:\ngot: %q\nwant match: %q", msg, expectedRegexp)
	}
}

func BenchmarkClientResponseFormatBodyStringTruncated(b *testing.B) {
	body := bytes.Repeat([]byte("test"), clientResponseBodyMaxSizeString)
	cr := &ClientResponse{
		Response: &http.Response{
			Status:     "400 Bad Request",
			StatusCode: http.StatusBadRequest,
			Proto:      "HTTP/1.1",
			Header: http.Header{
				"Foo": {"bar"},
			},
			ContentLength: int64(len(body)),
		},
		Body: body,
	}
	err := errors.New("error")
	err = WithClientResponse(err, cr)
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(io.Discard, "%+v", err)
	}
}

func TestClientResponseFormatBodyBytes(t *testing.T) {
	body := []byte{0xff}
	cr := &ClientResponse{
		Response: &http.Response{
			Status:     "400 Bad Request",
			StatusCode: http.StatusBadRequest,
			Proto:      "HTTP/1.1",
			Header: http.Header{
				"Foo": {"bar"},
			},
			ContentLength: int64(len(body)),
		},
		Body: body,
	}
	err := errors.New("error")
	err = WithClientResponse(err, cr)
	msg := fmt.Sprintf("%+v", err)
	expectedRegexp := regexp.MustCompile(`^HTTP client response
	status: 400 Bad Request
	proto: HTTP\/1\.1
	content length: 1
	headers:
		Foo: bar
	body:
00000000  ff                                                \|\.\|

stack
(\t.+\n)+error$`)
	if !expectedRegexp.MatchString(msg) {
		t.Fatalf("unexpected message:\ngot: %q\nwant match: %q", msg, expectedRegexp)
	}
}

func BenchmarkClientResponseFormatBodyBytes(b *testing.B) {
	body := []byte{0xff}
	cr := &ClientResponse{
		Response: &http.Response{
			Status:     "400 Bad Request",
			StatusCode: http.StatusBadRequest,
			Proto:      "HTTP/1.1",
			Header: http.Header{
				"Foo": {"bar"},
			},
			ContentLength: int64(len(body)),
		},
		Body: body,
	}
	err := errors.New("error")
	err = WithClientResponse(err, cr)
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(io.Discard, "%+v", err)
	}
}

func TestClientResponseFormatBodyBytesTruncated(t *testing.T) {
	body := bytes.Repeat([]byte{0xff}, clientResponseBodyMaxSizeBytes*2)
	cr := &ClientResponse{
		Response: &http.Response{
			Status:     "400 Bad Request",
			StatusCode: http.StatusBadRequest,
			Proto:      "HTTP/1.1",
			Header: http.Header{
				"Foo": {"bar"},
			},
			ContentLength: int64(len(body)),
		},
		Body: body,
	}
	err := errors.New("error")
	err = WithClientResponse(err, cr)
	msg := fmt.Sprintf("%+v", err)
	expectedRegexp := regexp.MustCompile(`^HTTP client response
	status: 400 Bad Request
	proto: HTTP\/1\.1
	content length: 2048
	headers:
		Foo: bar
	body:
(.{8}  ((ff ){8} ){2}\|\.{16}\|
)+\(truncated to 1024 bytes\)

stack
(\t.+\n)+error$`)
	if !expectedRegexp.MatchString(msg) {
		t.Fatalf("unexpected message:\ngot: %q\nwant match: %q", msg, expectedRegexp)
	}
}

func BenchmarkClientResponseFormatBodyBytesTruncated(b *testing.B) {
	body := bytes.Repeat([]byte{0xff}, clientResponseBodyMaxSizeBytes*2)
	cr := &ClientResponse{
		Response: &http.Response{
			Status:     "400 Bad Request",
			StatusCode: http.StatusBadRequest,
			Proto:      "HTTP/1.1",
			Header: http.Header{
				"Foo": {"bar"},
			},
			ContentLength: int64(len(body)),
		},
		Body: body,
	}
	err := errors.New("error")
	err = WithClientResponse(err, cr)
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(io.Discard, "%+v", err)
	}
}

func TestClientResponseFormatNoBody(t *testing.T) {
	cr := &ClientResponse{
		Response: &http.Response{
			Status:     "400 Bad Request",
			StatusCode: http.StatusBadRequest,
			Proto:      "HTTP/1.1",
			Header: http.Header{
				"Foo": {"bar"},
			},
		},
	}
	err := errors.New("error")
	err = WithClientResponse(err, cr)
	msg := fmt.Sprintf("%+v", err)
	expectedRegexp := regexp.MustCompile(`^HTTP client response
	status: 400 Bad Request
	proto: HTTP/1\.1
	headers:
		Foo: bar\n
stack
(\t.+\n)+error$`)
	if !expectedRegexp.MatchString(msg) {
		t.Fatalf("unexpected message:\ngot: %q\nwant match: %q", msg, expectedRegexp)
	}
}

func BenchmarkClientResponseFormatNoBody(b *testing.B) {
	cr := &ClientResponse{
		Response: &http.Response{
			Status:     "400 Bad Request",
			StatusCode: http.StatusBadRequest,
			Proto:      "HTTP/1.1",
			Header: http.Header{
				"Foo": {"bar"},
			},
		},
	}
	err := errors.New("error")
	err = WithClientResponse(err, cr)
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(io.Discard, "%+v", err)
	}
}
