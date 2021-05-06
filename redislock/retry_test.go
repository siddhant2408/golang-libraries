package redislock

import (
	"testing"
	"time"
)

func TestNoRetry(t *testing.T) {
	r := NoRetry()
	_, ok := r.Retry()
	if ok {
		t.Fatal(true)
	}
}

func TestLimitRetry(t *testing.T) {
	r := LimitRetry(
		&testRetry{
			d:  1 * time.Second,
			ok: true,
		},
		1,
	)
	d, ok := r.Retry()
	if !ok {
		t.Fatal(false)
	}
	if d != 1*time.Second {
		t.Fatalf("unexpected duration: got %v, want %v", d, 1*time.Second)
	}
	_, ok = r.Retry()
	if ok {
		t.Fatal(true)
	}
}

func TestDelayRetry(t *testing.T) {
	r := DelayRetry(1 * time.Second)
	d, ok := r.Retry()
	if !ok {
		t.Fatal(false)
	}
	if d != 1*time.Second {
		t.Fatalf("unexpected duration: got %v, want %v", d, 1*time.Second)
	}
}

func TestExponentialRetry(t *testing.T) {
	r := ExponentialRetry(1*time.Millisecond, 1*time.Second)
	for _, expectedD := range []time.Duration{
		1 * time.Millisecond,
		2 * time.Millisecond,
		4 * time.Millisecond,
		8 * time.Millisecond,
		16 * time.Millisecond,
		32 * time.Millisecond,
		64 * time.Millisecond,
		128 * time.Millisecond,
		256 * time.Millisecond,
		512 * time.Millisecond,
		1 * time.Second,
		1 * time.Second,
		1 * time.Second,
		1 * time.Second,
		1 * time.Second,
	} {
		d, ok := r.Retry()
		if !ok {
			t.Fatal(false)
		}
		if d != expectedD {
			t.Fatalf("unexpected duration: got %v, want %v", d, expectedD)
		}
	}
}

type testRetry struct {
	d  time.Duration
	ok bool
}

func (r *testRetry) Retry() (time.Duration, bool) {
	return r.d, r.ok
}
