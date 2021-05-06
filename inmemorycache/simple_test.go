package inmemorycache

import (
	"testing"
)

func TestSimpleSetGet(t *testing.T) {
	c := newTestSimple()
	c.Set(testKey, testValue)
	vi, found := c.Get(testKey)
	if !found {
		t.Fatal("not found")
	}
	if vi != testValue {
		t.Fatalf("unexpected value: got %v, want %q", vi, testValue)
	}
}

func TestSimpleGetNotFound(t *testing.T) {
	c := newTestSimple()
	_, found := c.Get(testKey)
	if found {
		t.Fatal("found")
	}
}

func TestSimpleRemove(t *testing.T) {
	c := newTestSimple()
	c.Set(testKey, testValue)
	c.Remove(testKey)
	_, found := c.Get(testKey)
	if found {
		t.Fatal("found")
	}
}

func BenchmarkSimpleSetGet(b *testing.B) {
	c := newTestSimple()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(testKey, testValue)
		c.Get(testKey)
	}
}

func BenchmarkSimpleGet(b *testing.B) {
	c := newTestSimple()
	c.Set(testKey, testValue)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(testKey)
	}
}

func newTestSimple() Cache {
	return New()
}
