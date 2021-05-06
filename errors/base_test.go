package errors_test

import (
	"io"
	"testing"

	. "github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/errors/internal"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestNew(t *testing.T) {
	err := New("error")
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
	sfs := StackFrames(err)
	if len(sfs) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(sfs), 1)
	}
}

func TestNewf(t *testing.T) {
	err := Newf("%s", "error")
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
	sfs := StackFrames(err)
	if len(sfs) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(sfs), 1)
	}
}

func TestIs(t *testing.T) {
	err := io.EOF
	err = Wrap(err, "test")
	ok := Is(err, io.EOF)
	if !ok {
		t.Fatal("not ok")
	}
}

func TestUnwrapAll(t *testing.T) {
	rootErr := internal.NewBase("error")
	err := WithMessage(rootErr, "test1")
	err = WithMessage(err, "test2")
	err = WithMessage(err, "test3")
	err = WithMessage(err, "test4")
	err = WithMessage(err, "test5")
	werr := UnwrapAll(err)
	if werr != rootErr { //nolint:goerr113 // We want to compare the current error.
		t.Fatalf("unexpected wrapped: got %q, want %q", werr, rootErr)
	}
}

func TestUnwrapAllNone(t *testing.T) {
	err := internal.NewBase("error")
	werr := UnwrapAll(err)
	if werr != err { //nolint:goerr113 // We want to compare the current error.
		t.Fatalf("unexpected wrapped: got %q, want %q", werr, err)
	}
}

func TestUnwrapAllNil(t *testing.T) {
	err := UnwrapAll(nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
