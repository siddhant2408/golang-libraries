package httpclientrequest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func Test(t *testing.T) {
	ctx := context.Background()
	body := []byte("test")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write(body)
	}))
	defer srv.Close()
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	res, err := Do(ctx, req)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if res.Response == nil {
		t.Fatal("response is nil")
	}
	if !bytes.Equal(res.Body, body) {
		t.Fatalf("unexpected response body: got %v, want %v", res.Body, body)
	}
}

func TestClient(t *testing.T) {
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {}))
	defer srv.Close()
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	clt := new(http.Client)
	_, err = Do(ctx, req, Client(clt))
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestCopyBody(t *testing.T) {
	ctx := context.Background()
	body := []byte("test")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write(body)
	}))
	defer srv.Close()
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	buf := new(bytes.Buffer)
	res, err := Do(ctx, req, CopyBody(buf))
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if res.Response == nil {
		t.Fatal("response is nil")
	}
	if res.Body != nil {
		t.Fatal("body not nil")
	}
	if !bytes.Equal(buf.Bytes(), body) {
		t.Fatalf("unexpected body copy: got %v, want %v", buf.Bytes(), body)
	}
}

func TestStatus(t *testing.T) {
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {}))
	defer srv.Close()
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	_, err = Do(ctx, req, Status(http.StatusOK))
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestErrorTimeout(t *testing.T) {
	ctx := context.Background()
	interrupt := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		<-interrupt
	}))
	defer srv.Close()
	defer close(interrupt)
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	res, err := Do(ctx, req, Timeout(1*time.Millisecond))
	if err == nil {
		t.Fatal("no error")
	}
	if res.Response != nil {
		t.Fatal("response not nil")
	}
	if len(res.Body) != 0 {
		t.Fatal("response body not empty")
	}
}

func TestErrorMaxBodySize(t *testing.T) {
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte("very long response"))
	}))
	defer srv.Close()
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	res, err := Do(ctx, req, MaxBodySize(8))
	if err == nil {
		t.Fatal("no error")
	}
	if res.Response == nil {
		t.Fatal("response is nil")
	}
	expectedBody := []byte("very lon")
	if !bytes.Equal(res.Body, expectedBody) {
		t.Fatalf("unexpected response body: got %v, want %v", res.Body, expectedBody)
	}
}

func TestErrorStatus2XX(t *testing.T) {
	ctx := context.Background()
	body := []byte("bad request")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(body)
	}))
	defer srv.Close()
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	res, err := Do(ctx, req, Status2XX())
	if err == nil {
		t.Fatal("no error")
	}
	if res.Response == nil {
		t.Fatal("response is nil")
	}
	if !bytes.Equal(res.Body, body) {
		t.Fatalf("unexpected response body: got %v, want %v", res.Body, body)
	}
}

func TestErrorStatus(t *testing.T) {
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	res, err := Do(ctx, req, Status(http.StatusOK))
	if err == nil {
		t.Fatal("no error")
	}
	if res.Response == nil {
		t.Fatal("response is nil")
	}
}

func TestPanicResponseBodyRead(t *testing.T) {
	ctx := context.Background()
	body := []byte("test")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write(body)
	}))
	defer srv.Close()
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	res, err := Do(ctx, req)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatal("no panic")
		}
	}()
	_, _ = res.Response.Body.Read(nil)
}

func TestPanicResponseBodyClose(t *testing.T) {
	ctx := context.Background()
	body := []byte("test")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write(body)
	}))
	defer srv.Close()
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	res, err := Do(ctx, req)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatal("no panic")
		}
	}()
	_ = res.Response.Body.Close()
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Errorf("unexpected method: got %s, want %s", req.Method, http.MethodGet)
		}
	}))
	defer srv.Close()
	_, err := Get(ctx, srv.URL)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestGetErrorRequest(t *testing.T) {
	ctx := context.Background()
	_, err := Get(ctx, "invalid\tinvalid")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestHead(t *testing.T) {
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodHead {
			t.Errorf("unexpected method: got %s, want %s", req.Method, http.MethodHead)
		}
	}))
	defer srv.Close()
	_, err := Head(ctx, srv.URL)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestPost(t *testing.T) {
	ctx := context.Background()
	body := []byte("test")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			t.Errorf("unexpected method: got %s, want %s", req.Method, http.MethodPost)
		}
		b, err := io.ReadAll(req.Body)
		if err != nil {
			testutils.ErrorErr(t, err)
			return
		}
		if !bytes.Equal(b, body) {
			t.Errorf("body not equal: got %v, want %v", b, body)
		}
	}))
	defer srv.Close()
	_, err := Post(ctx, srv.URL, "text/plain", bytes.NewReader(body))
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestPostForm(t *testing.T) {
	ctx := context.Background()
	data := url.Values{
		"foo": []string{"bar"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			t.Errorf("unexpected method: got %s, want %s", req.Method, http.MethodPost)
		}
		b, err := io.ReadAll(req.Body)
		if err != nil {
			testutils.ErrorErr(t, err)
			return
		}
		q, err := url.ParseQuery(string(b))
		if err != nil {
			testutils.ErrorErr(t, err)
			return
		}
		testutils.CompareError(t, "unexpected form", q, data)
	}))
	defer srv.Close()
	_, err := PostForm(ctx, srv.URL, data)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
