package inmemorycache

import (
	"testing"
)

func TestLRUSetGet(t *testing.T) {
	c := newTestLRU()
	c.Set(testKey, testValue)
	vi, found := c.Get(testKey)
	if !found {
		t.Fatal("not found")
	}
	if vi != testValue {
		t.Fatalf("unexpected value: got %v, want %q", vi, testValue)
	}
}

func TestLRUGetNotFound(t *testing.T) {
	c := newTestLRU()
	_, found := c.Get(testKey)
	if found {
		t.Fatal("found")
	}
}

func TestLRURemove(t *testing.T) {
	c := newTestLRU()
	c.Set(testKey, testValue)
	c.Remove(testKey)
	_, found := c.Get(testKey)
	if found {
		t.Fatal("found")
	}
}

func TestLRUSetEvict(t *testing.T) {
	c := newTestLRU()
	c.Set(testKey, testValue)
	for i := 0; i < testMaxSize; i++ {
		c.Set(i, testValue)
	}
	_, found := c.Get(testKey)
	if found {
		t.Fatal("found")
	}
}

func TestLRUSetOverwrite(t *testing.T) {
	c := newTestLRU()
	c.Set(testKey, "test")
	c.Set(testKey, testValue)
	vi, found := c.Get(testKey)
	if !found {
		t.Fatal("not found")
	}
	if vi != testValue {
		t.Fatalf("unexpected value: got %v, want %q", vi, testValue)
	}
}

func BenchmarkLRUSetGet(b *testing.B) {
	c := newTestLRU()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(testKey, testValue)
		c.Get(testKey)
	}
}

func BenchmarkLRUGet(b *testing.B) {
	c := newTestLRU()
	c.Set(testKey, testValue)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(testKey)
	}
}

func BenchmarkLRUSetEvict(b *testing.B) {
	c := newTestLRU()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(i, testValue)
	}
}

func BenchmarkLRUSetGetEvict(b *testing.B) {
	c := newTestLRU()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(i, testValue)
		c.Get(i)
	}
}

func newTestLRU() Cache {
	return New(
		MaxSize(testMaxSize),
	)
}
