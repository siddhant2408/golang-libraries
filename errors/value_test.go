package errors_test

import (
	"fmt"
	"io"
	"testing"

	. "github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/errors/internal"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestValue(t *testing.T) {
	err := internal.NewBase("error")
	err = WithValue(err, "foo", "bar")
	vals := Values(err)
	expected := map[string]interface{}{
		"foo": "bar",
	}
	testutils.Compare(t, "unexpected values", vals, expected)
}

func TestValueOverWrite(t *testing.T) {
	err := internal.NewBase("error")
	err = WithValue(err, "test", 1)
	err = WithValue(err, "test", 2)
	vals := Values(err)
	expected := map[string]interface{}{
		"test": 2,
	}
	testutils.Compare(t, "unexpected values", vals, expected)
}

func TestValueNil(t *testing.T) {
	err := WithValue(nil, "foo", "bar")
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestValuesEmpty(t *testing.T) {
	err := internal.NewBase("error")
	vals := Values(err)
	if len(vals) != 0 {
		t.Fatalf("values not empty: got %#v", vals)
	}
}

func TestValueError(t *testing.T) {
	err := internal.NewBase("error")
	err = WithValue(err, "foo", "bar")
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestValueFormat(t *testing.T) {
	err := internal.NewBase("error")
	err = WithValue(err, "foo", "bar")
	s := fmt.Sprintf("%+v", err)
	expected := "value foo = (string) (len=3) \"bar\"\nerror"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func BenchmarkValueFormat(b *testing.B) {
	err := internal.NewBase("error")
	err = WithValue(err, "foo", "bar")
	for i := 0; i < b.N; i++ {
		_, _ = fmt.Fprintf(io.Discard, "%+v", err)
	}
}
