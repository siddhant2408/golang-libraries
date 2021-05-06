package errors_test

import (
	"fmt"
	"io"
	"testing"

	. "github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/errors/internal"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestIgnore(t *testing.T) {
	err := internal.NewBase("error")
	err = Ignore(err)
	ignored := IsIgnored(err)
	if !ignored {
		t.Fatalf("unexpected ignored: got %t, want %t", ignored, true)
	}
}

func TestIgnoreNil(t *testing.T) {
	err := Ignore(nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestIsIgnoredFalse(t *testing.T) {
	err := internal.NewBase("error")
	ignored := IsIgnored(err)
	if ignored {
		t.Fatalf("unexpected ignored: got %t, want %t", ignored, false)
	}
}

func TestIgnoreError(t *testing.T) {
	err := internal.NewBase("error")
	err = Ignore(err)
	s := err.Error()
	expected := "ignored: error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestIgnoreFormat(t *testing.T) {
	err := internal.NewBase("error")
	err = Ignore(err)
	s := fmt.Sprint(err)
	expected := "ignored: error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func BenchmarkIgnoreFormat(b *testing.B) {
	err := internal.NewBase("error")
	err = Ignore(err)
	for i := 0; i < b.N; i++ {
		_, _ = fmt.Fprint(io.Discard, err)
	}
}
