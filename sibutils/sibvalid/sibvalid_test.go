package sibvalid

import (
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func testStrings(t *testing.T, f func(testing.TB, string), ss []string) {
	for _, s := range ss {
		t.Run(s, func(t *testing.T) {
			f(t, s)
		})
	}
}

func benchmarkStrings(b *testing.B, f func(testing.TB, string), ss []string) {
	for _, s := range ss {
		b.Run(s, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				f(b, s)
			}
		})
	}
}

func testStringValid(tb testing.TB, f func(string) error, s string) {
	err := f(s)
	if err != nil {
		testutils.FatalErr(tb, err)
	}
}

func testStringInvalid(tb testing.TB, f func(string) error, s string) {
	err := f(s)
	if err == nil {
		tb.Fatal("no error")
	}
}
