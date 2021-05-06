package httpjson

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/jsonreader"
	"github.com/siddhant2408/golang-libraries/jsontest"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestDecodeRequestBody(t *testing.T) {
	ctx := context.Background()
	data := map[string]interface{}{
		"foo": "bar",
	}
	body := jsonreader.New(data, nil)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", body)
	var res interface{}
	err := DecodeRequestBody(ctx, req, &res,
		RequestDecoder(func(dec *json.Decoder) {}),
	)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	testutils.Compare(t, "unexpected result", res, data)
}

func TestDecodeRequestBodyRemainingDataEmpty(t *testing.T) {
	ctx := context.Background()
	body := strings.NewReader("{} \n	\n ")
	req := httptest.NewRequest(http.MethodGet, "http://localhost", body)
	var res interface{}
	err := DecodeRequestBody(ctx, req, &res)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestDecodeRequestBodyErrorRead(t *testing.T) {
	ctx := context.Background()
	body := &errReader{
		err: errors.New("error"),
	}
	req := httptest.NewRequest(http.MethodGet, "http://localhost", body)
	var res interface{}
	err := DecodeRequestBody(ctx, req, &res)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestDecodeRequestBodyErrorUnmarshal(t *testing.T) {
	ctx := context.Background()
	body := strings.NewReader("invalid")
	req := httptest.NewRequest(http.MethodGet, "http://localhost", body)
	var res interface{}
	err := DecodeRequestBody(ctx, req, &res)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestDecodeRequestBodyErrorRemainingDataToken(t *testing.T) {
	ctx := context.Background()
	body := strings.NewReader("{}{}")
	req := httptest.NewRequest(http.MethodGet, "http://localhost", body)
	var res interface{}
	err := DecodeRequestBody(ctx, req, &res)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestDecodeRequestBodyErrorRemainingDataInvalid(t *testing.T) {
	ctx := context.Background()
	body := strings.NewReader("{}invalid")
	req := httptest.NewRequest(http.MethodGet, "http://localhost", body)
	var res interface{}
	err := DecodeRequestBody(ctx, req, &res)
	if err == nil {
		t.Fatal("no error")
	}
}

func BenchmarkDecodeRequestBody(b *testing.B) {
	ctx := context.Background()
	data := map[string]interface{}{
		"foo": "bar",
	}
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	for i := 0; i < b.N; i++ {
		body := jsonreader.New(data, nil)
		req.Body = io.NopCloser(body)
		var res interface{}
		err := DecodeRequestBody(ctx, req, &res)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
}

func TestWriteResponse(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()
	code := http.StatusOK
	data := map[string]interface{}{
		"foo": "bar",
	}
	err := WriteResponse(ctx, w, code, data,
		ResponseEncoder(func(enc *json.Encoder) {}),
		ResponseHeader(func(hd http.Header) {
			hd.Set("foo", "bar")
		}),
	)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	w.Flush()
	if w.Code != code {
		t.Fatalf("unexpected code: got %d, want %d", w.Code, code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected content type: got %q, want %q", w.Header().Get("Content-Type"), "application/json")
	}
	var res interface{}
	jsontest.Unmarshal(t, w.Body.Bytes(), &res)
	testutils.Compare(t, "unexpected result", res, data)
}

func BenchmarkWriteResponse(b *testing.B) {
	ctx := context.Background()
	code := http.StatusOK
	data := map[string]interface{}{
		"foo": "bar",
	}
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		err := WriteResponse(ctx, w, code, data)
		if err != nil {
			testutils.FatalErr(b, err)
		}
		w.Flush()
	}
}

type errReader struct {
	err error
}

func (r *errReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}
