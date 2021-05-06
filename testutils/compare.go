package testutils

import (
	"reflect" //nolint:depguard // Required for type comparison.
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pierrre/compare"
)

// Compare is an alias for CompareFatal.
// It is NOT deprecated.
func Compare(tb testing.TB, msg string, got, want interface{}) {
	tb.Helper()
	CompareFatal(tb, msg, got, want)
}

// CompareFatal compares 2 values and fails immediately the test if there is a difference.
func CompareFatal(tb testing.TB, msg string, got, want interface{}) {
	tb.Helper()
	comparef(tb, msg, got, want, tb.Fatalf)
}

// CompareError compares 2 values and flags the test as errored if there is a difference.
// It is useful for tests that run in a separate goroutine.
// You can get the "failed" status of the current test with tb.Failed().
func CompareError(tb testing.TB, msg string, got, want interface{}) {
	tb.Helper()
	comparef(tb, msg, got, want, tb.Errorf)
}

func comparef(tb testing.TB, msg string, got, want interface{}, f func(format string, args ...interface{})) {
	tb.Helper()
	checkCompareComparable(tb, got, want, f)
	if tb.Failed() {
		return
	}
	diff := compare.Compare(got, want)
	if len(diff) != 0 {
		f("%s:\ngot:\n%s\nwant:\n%s\ndiff:\n%+v", msg, spew.Sdump(got), spew.Sdump(want), diff)
	}
}

func checkCompareComparable(tb testing.TB, got, want interface{}, f func(format string, args ...interface{})) {
	tb.Helper()
	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		return
	}
	switch got.(type) {
	case string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, complex64, complex128, error:
		f("%T values must be compared with == or !=", got)
	case []byte:
		f("%T values must be compared with bytes.Equal()", got)
	}
}
