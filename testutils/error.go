package testutils

import (
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
)

// FatalErr calls tb.Fatal() with an improved error formatting.
func FatalErr(tb testing.TB, err error) {
	tb.Helper()
	err = errors.Wrap(err, "")
	tb.Fatalf("%v\n%+v", err, err)
}

// ErrorErr calls tb.Error() with an improved error formatting.
func ErrorErr(tb testing.TB, err error) {
	tb.Helper()
	err = errors.Wrap(err, "")
	tb.Errorf("%v\n%+v", err, err)
}

// SkipErr calls tb.Skip() with an improved error formatting.
func SkipErr(tb testing.TB, err error) {
	tb.Helper()
	err = errors.Wrap(err, "")
	tb.Skipf("%v\n%+v", err, err)
}

// LogErr calls tb.Log() with an improved error formatting.
func LogErr(tb testing.TB, err error) {
	tb.Helper()
	err = errors.Wrap(err, "")
	tb.Logf("%v\n%+v", err, err)
}
