package jsonlog

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func Test(t *testing.T) {
	ctx := context.Background()
	buf := new(bytes.Buffer)
	l := New(buf)
	err := l.Log(ctx, &testData{
		Test: "test1",
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = l.Log(ctx, &testData{
		Test: "test2",
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	res := buf.String()
	expected := `{"test":"test1"}` + "\n" + `{"test":"test2"}` + "\n"
	if res != expected {
		t.Fatalf("unexpected result: got %q, want %q", res, expected)
	}
}

func TestErrorEncode(t *testing.T) {
	l := New(io.Discard)
	err := l.Log(context.Background(), func() {})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestFile(t *testing.T) {
	d, cleanD := testutils.TempDir(t, "", "")
	defer cleanD()
	name := filepath.Join(d, "test.log")
	l, closeL, err := NewFile(name, os.FileMode(0644))
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = l.Log(context.Background(), &testData{
		Test: "test",
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = closeL()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	content, err := os.ReadFile(name)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	res := string(content)
	expected := `{"test":"test"}` + "\n"
	if res != expected {
		t.Fatalf("unexpected content: got %q, want %q", res, expected)
	}
}

func TestFileError(t *testing.T) {
	_, _, err := NewFile("", os.FileMode(0644))
	if err == nil {
		t.Fatal("no error")
	}
}

func Benchmark(b *testing.B) {
	ctx := context.Background()
	l := New(io.Discard)
	data := &testData{
		Test: "test",
	}
	for i := 0; i < b.N; i++ {
		err := l.Log(ctx, data)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
}

func BenchmarkFile(b *testing.B) {
	ctx := context.Background()
	data := &testData{
		Test: "test",
	}
	d, cleanD := testutils.TempDir(b, "", "")
	defer cleanD()
	name := filepath.Join(d, "test.log")
	l, closeL, err := NewFile(name, os.FileMode(0644))
	if err != nil {
		testutils.FatalErr(b, err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = l.Log(ctx, data)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
	err = closeL()
	if err != nil {
		testutils.FatalErr(b, err)
	}
}

func TestOptionalFile(t *testing.T) {
	d, cleanD := testutils.TempDir(t, "", "")
	defer cleanD()
	name := filepath.Join(d, "test.log")
	l, closeL, err := NewOptionalFile(name, os.FileMode(0644))
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = l.Log(context.Background(), &testData{
		Test: "test",
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = closeL()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	content, err := os.ReadFile(name)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	res := string(content)
	expected := `{"test":"test"}` + "\n"
	if res != expected {
		t.Fatalf("unexpected content: got %q, want %q", res, expected)
	}
}

func TestOptionalFileEmpty(t *testing.T) {
	l, closeL, err := NewOptionalFile("", os.FileMode(0644))
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = l.Log(context.Background(), &testData{
		Test: "test",
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = closeL()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestError(t *testing.T) {
	tl := &testLogger{}
	el := &Error{
		Logger: tl,
	}
	el.Log(context.Background(), &testData{
		Test: "test",
	})
}

func TestErrorError(t *testing.T) {
	tl := &testLogger{
		err: errors.New("error"),
	}
	var called testutils.CallCounter
	el := &Error{
		Logger: tl,
		OnError: func(_ context.Context, err error) {
			called.Call()
			if err == nil {
				t.Fatal("no error")
			}
		},
	}
	el.Log(context.Background(), &testData{
		Test: "test",
	})
	called.AssertCalled(t)
}

func TestNewErrorOptionalFile(t *testing.T) {
	d, cleanD := testutils.TempDir(t, "", "")
	defer cleanD()
	name := filepath.Join(d, "test.log")
	l, closeL, err := NewErrorOptionalFile(name, os.FileMode(0644), func(ctx context.Context, err error) {})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if l == nil {
		t.Fatal("no logger")
	}
	err = closeL()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

type testLogger struct {
	callCounter testutils.CallCounter
	err         error
}

func (l *testLogger) Log(ctx context.Context, data interface{}) error {
	l.callCounter.Call()
	return l.err
}

type testData struct {
	Test string `json:"test"`
}
