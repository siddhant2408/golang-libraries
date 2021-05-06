package sibhttpua

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/siddhant2408/golang-libraries/httpclientrequest"
	"github.com/siddhant2408/golang-libraries/testutils"
)

const (
	testAppName  = "test"
	testVersion  = "1.2.3"
	testExpected = "Siddhant/1.0 (test 1.2.3; +https://siddhant-test.com)"
)

func TestGet(t *testing.T) {
	ua := Get(testAppName, testVersion)
	if ua != testExpected {
		t.Fatalf("unexpected user-agent: got %q, want %q", ua, testExpected)
	}
}

func TestGetPanicCheckAppName(t *testing.T) {
	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatal("no panic")
		}
	}()
	Get("", testVersion)
}

func TestGetPanicCheckVersion(t *testing.T) {
	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatal("no panic")
		}
	}()
	Get(testAppName, "")
}

func TestSetToRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	SetToRequest(req, testAppName, testVersion)
	ua := req.UserAgent()
	if ua != testExpected {
		t.Fatalf("unexpected user-agent: got %q, want %q", ua, testExpected)
	}
}

func TestAlreadySet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	customUA := "foobar"
	req.Header.Set("User-Agent", customUA)
	SetToRequest(req, testAppName, testVersion)
	ua := req.UserAgent()
	if ua != customUA {
		t.Fatalf("unexpected user-agent: got %q, want %q", ua, customUA)
	}
}

func TestWrapRoundTripper(t *testing.T) {
	// Reset to the original value at the end of the test.
	dt := http.DefaultTransport
	defer func() {
		http.DefaultTransport = dt
	}()
	WrapDefaultTransport(testAppName, testVersion)
	h := http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		ua := req.UserAgent()
		if ua != testExpected {
			t.Errorf("unexpected user-agent: got %q, want %q", ua, testExpected)
		}
	})
	srv := httptest.NewServer(h)
	defer srv.Close()
	ctx := context.Background()
	_, err := httpclientrequest.Get(ctx, srv.URL)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
