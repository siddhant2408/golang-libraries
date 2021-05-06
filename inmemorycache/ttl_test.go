package inmemorycache

import (
	"testing"
	"time"

	"github.com/siddhant2408/golang-libraries/timeutils"
)

const testTTL = 1 * time.Minute

func TestTTLSetGet(t *testing.T) {
	c := newTestTTL()
	c.Set(testKey, testValue)
	vi, found := c.Get(testKey)
	if !found {
		t.Fatal("not found")
	}
	if vi != testValue {
		t.Fatalf("unexpected value: got %v, want %q", vi, testValue)
	}
}

func TestTTLGetNotFound(t *testing.T) {
	c := newTestTTL()
	_, found := c.Get(testKey)
	if found {
		t.Fatal("found")
	}
}

func TestTTLGetExpired(t *testing.T) {
	c := newTestTTL()
	c.Set(testKey, testValue)
	// Move forward the fake time by the test TTL + 1 second, so the key expires.
	timeutils.SetFixed(timeutils.Now().Add(testTTL + 1*time.Second))
	_, found := c.Get(testKey)
	if found {
		t.Fatal("found")
	}
}

func BenchmarkTTLSetGet(b *testing.B) {
	c := newTestTTL()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(testKey, testValue)
		c.Get(testKey)
	}
}

func BenchmarkTTLGet(b *testing.B) {
	c := newTestTTL()
	c.Set(testKey, testValue)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(testKey)
	}
}

func newTestTTL() Cache {
	return New(
		TTL(testTTL),
	)
}
