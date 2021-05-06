package errors_test

import (
	"fmt"
	"io"
	"testing"

	. "github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/errors/internal"
)

func TestFormat(t *testing.T) {
	err := internal.NewBase("error")
	err = withTestFormat(err)
	for _, tc := range []struct {
		ft       string
		expected string
	}{
		{
			ft:       "%v",
			expected: "test: error",
		},
		{
			ft:       "%+v",
			expected: "test verbose\nerror",
		},
		{
			ft:       "%s",
			expected: "test: error",
		},
		{
			ft:       "%q",
			expected: "\"test: error\"",
		},
	} {
		t.Run(tc.ft, func(t *testing.T) {
			s := fmt.Sprintf(tc.ft, err)
			if s != tc.expected {
				t.Fatalf("unexpected message: got %q, want %q", s, tc.expected)
			}
		})
	}
}

func BenchmarkFormat(b *testing.B) {
	for wrapCount := 1; wrapCount <= 64; wrapCount *= 2 {
		b.Run(fmt.Sprintf("WrapCount_%d", wrapCount), func(b *testing.B) {
			err := internal.NewBase("error")
			for i := 0; i < wrapCount; i++ {
				err = withTestFormat(err)
			}
			for _, ft := range []string{"%v", "%+v", "%q"} {
				b.Run(ft, func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						_, _ = fmt.Fprintf(io.Discard, ft, err)
					}
				})
			}
		})
	}
}

type testFormat struct {
	err error
}

func withTestFormat(err error) error {
	if err == nil {
		return nil
	}
	return &testFormat{
		err: err,
	}
}

func (err *testFormat) WriteErrorMessage(w Writer, verbose bool) bool {
	if verbose {
		_, _ = w.WriteString("test verbose")
	} else {
		_, _ = w.WriteString("test")
	}
	return true
}

func (err *testFormat) Error() string                 { return Error(err) }
func (err *testFormat) Format(s fmt.State, verb rune) { Format(err, s, verb) }
func (err *testFormat) Unwrap() error                 { return err.err }
