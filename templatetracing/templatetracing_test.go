package templatetracing

import (
	"bytes"
	"context"
	"testing"
	"text/template"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	s := "aaa {{ .test }} bbb"
	tmpl, err := template.New("test").Parse(s)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	buf := new(bytes.Buffer)
	data := map[string]interface{}{
		"test": "zzz",
	}
	err = Execute(ctx, buf, tmpl, data)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	res := buf.String()
	expected := "aaa zzz bbb"
	if res != expected {
		t.Fatalf("unexpected result: got %q, want %q", res, expected)
	}
}

func TestExecuteError(t *testing.T) {
	ctx := context.Background()
	s := "aaa {{ .test }} bbb"
	tmpl, err := template.New("test").Parse(s)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	buf := new(bytes.Buffer)
	data := struct{}{}
	err = Execute(ctx, buf, tmpl, data)
	if err == nil {
		t.Fatal("no error")
	}
}
