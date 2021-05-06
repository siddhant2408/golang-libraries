package httptestlocal_test

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/httptestlocal"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestNormal(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", bytes.NewReader(nil))
	rt := httptestlocal.WrapRoundTripper(&testRoundTripper{
		resp: &http.Response{
			Body: io.NopCloser(bytes.NewReader(nil)),
		},
	})
	resp, err := rt.RoundTrip(req)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	_ = resp.Body.Close()
}

func TestContextAllow(t *testing.T) {
	ctx := context.Background()
	ctx = httptestlocal.Allow(ctx)
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	req = req.WithContext(ctx)
	rt := httptestlocal.WrapRoundTripper(&testRoundTripper{
		resp: &http.Response{
			Body: io.NopCloser(bytes.NewReader(nil)),
		},
	})
	resp, err := rt.RoundTrip(req)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	_ = resp.Body.Close()
}

func TestErrorExternal(t *testing.T) {
	_, err := net.LookupIP("example.com")
	if err != nil {
		err = errors.Wrap(err, "test requires to resolve DNS")
		testutils.SkipErr(t, err)
	}
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	rt := httptestlocal.WrapRoundTripper(&testRoundTripper{
		resp: &http.Response{
			Body: io.NopCloser(bytes.NewReader(nil)),
		},
	})
	_, err = rt.RoundTrip(req) //nolint:bodyclose // An error should be returned, so it's not required to close the body.
	if err == nil {
		t.Fatal("no error")
	}
}

func TestErrorLookup(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://invalid", bytes.NewReader(nil))
	trt := &testRoundTripper{
		err: errors.New("error"),
	}
	rt := httptestlocal.WrapRoundTripper(trt)
	_, err := rt.RoundTrip(req) //nolint:bodyclose // An error should be returned, so it's not required to close the body.
	if err == nil {
		t.Fatal("no error")
	}
	trt.callCounter.AssertCalled(t)
}

type testRoundTripper struct {
	resp        *http.Response
	err         error
	callCounter testutils.CallCounter
}

func (rt *testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.callCounter.Call()
	if req.Body != nil {
		_ = req.Body.Close()
	}
	return rt.resp, rt.err
}
