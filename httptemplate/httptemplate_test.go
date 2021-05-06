package httptemplate

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestWriteResponse(t *testing.T) {
	ctx := context.Background()
	s := "aaa {{ .test }} bbb"
	tmpl, err := template.New("test").Parse(s)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	w := httptest.NewRecorder()
	hdf := func(hd http.Header) {
		hd.Set("foo", "bar")
	}
	code := http.StatusOK
	data := map[string]interface{}{
		"test": "zzz",
	}
	err = WriteResponse(ctx, w, hdf, code, tmpl, data)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	w.Flush()
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", w.Code, http.StatusOK)
	}
	res := w.Body.String()
	expected := "aaa zzz bbb"
	if res != expected {
		t.Fatalf("unexpected result: got %q, want %q", res, expected)
	}
}

func TestWriteResponseError(t *testing.T) {
	ctx := context.Background()
	s := "aaa {{ .test }} bbb"
	tmpl, err := template.New("test").Parse(s)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	w := httptest.NewRecorder()
	hdf := func(hd http.Header) {
		hd.Set("foo", "bar")
	}
	code := http.StatusOK
	data := struct{}{}
	err = WriteResponse(ctx, w, hdf, code, tmpl, data)
	if err == nil {
		t.Fatal("no error")
	}
}
