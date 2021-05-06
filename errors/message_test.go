package errors_test

import (
	"fmt"
	"io"
	"testing"

	. "github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/errors/internal"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestMessage(t *testing.T) {
	err := internal.NewBase("error")
	err = WithMessagef(err, "%s", "test")
	s := err.Error()
	expected := "test: error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestMessageNil(t *testing.T) {
	err := WithMessage(nil, "test")
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestMessageEmpty(t *testing.T) {
	err := internal.NewBase("error")
	err = WithMessage(err, "")
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestMessageFormat(t *testing.T) {
	err := internal.NewBase("error")
	err = WithMessagef(err, "%s", "test")
	s := fmt.Sprint(err)
	expected := "test: error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestWrap(t *testing.T) {
	err := internal.NewBase("error")
	err = Wrap(err, "test1")
	err = Wrapf(err, "%s", "test2")
	s := err.Error()
	expected := "test2: test1: error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
	sfs := StackFrames(err)
	if len(sfs) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(sfs), 1)
	}
}

func BenchmarkMessageFormat(b *testing.B) {
	err := internal.NewBase("error")
	err = WithMessage(err, "test")
	for i := 0; i < b.N; i++ {
		_, _ = fmt.Fprint(io.Discard, err)
	}
}
