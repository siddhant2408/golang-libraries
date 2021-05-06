package timeutils_test

import (
	"testing"
	"time"

	"github.com/siddhant2408/golang-libraries/testutils"
	. "github.com/siddhant2408/golang-libraries/timeutils"
)

func TestSetFixed(t *testing.T) {
	tm1 := time.Now()
	SetFixed(tm1)
	tm2 := Now()
	if !tm1.Equal(tm2) {
		t.Fatalf("unexpected result: got %s, want %s", tm2, tm1)
	}
}

func TestSince(t *testing.T) {
	InitFixed()
	d1 := 1 * time.Second
	tm := Now().Add(-d1)
	d2 := Since(tm)
	if d2 != d1 {
		t.Fatalf("unexpected duration: got %s, want %s", d2, d1)
	}
}

func TestUntil(t *testing.T) {
	InitFixed()
	d1 := 1 * time.Second
	tm := Now().Add(d1)
	d2 := Until(tm)
	if d2 != d1 {
		t.Fatalf("unexpected duration: got %s, want %s", d2, d1)
	}
}

const testLocationName = "Europe/Paris"

func TestLoadLocation(t *testing.T) {
	l, err := LoadLocation(testLocationName)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if l == nil {
		t.Fatal("nil")
	}
}

func TestLoadLocationAllocs(t *testing.T) {
	allocs := testing.AllocsPerRun(1000, func() {
		_, err := LoadLocation(testLocationName)
		if err != nil {
			testutils.FatalErr(t, err)
		}
	})
	if allocs != 0 {
		t.Fatalf("unexpected allocs: got %f, want %d", allocs, 0)
	}
}

func TestLoadLocationError(t *testing.T) {
	_, err := LoadLocation("invalid")
	if err == nil {
		t.Fatal("no error")
	}
}

func BenchmarkLoadLocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := LoadLocation(testLocationName)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
}

func BenchmarkLoadLocationStdlib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := time.LoadLocation(testLocationName)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
}
