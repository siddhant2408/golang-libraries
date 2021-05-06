package inmemorycache

import (
	"fmt"
	"sync"
	"testing"

	"github.com/siddhant2408/golang-libraries/goroutine"
)

func TestSyncMapSetGet(t *testing.T) {
	c := newTestSyncMap()
	wg := new(sync.WaitGroup)
	for g := 0; g < testConcurrentGoroutines; g++ {
		g := g
		goroutine.WaitGroup(wg, func() {
			k := fmt.Sprintf("test_key_%d", g)
			v1 := fmt.Sprintf("test_value_%d", g)
			for i := 0; i < testConcurrentIterations; i++ {
				c.Set(k, v1)
				v2i, found := c.Get(k)
				if !found {
					t.Error("not found")
					return
				}
				if v2i != v1 {
					t.Errorf("unexpected value: got %v, want %q", v2i, v1)
					return
				}
			}
		})
	}
	wg.Wait()
}

func TestSyncMapGet(t *testing.T) {
	c := newTestSyncMap()
	c.Set(testKey, testValue)
	wg := new(sync.WaitGroup)
	for g := 0; g < testConcurrentGoroutines; g++ {
		goroutine.WaitGroup(wg, func() {
			for i := 0; i < testConcurrentIterations; i++ {
				vi, found := c.Get(testKey)
				if !found {
					t.Error("not found")
					return
				}
				if vi != testValue {
					t.Errorf("unexpected value: got %v, want %q", vi, testValue)
					return
				}
			}
		})
	}
	wg.Wait()
}

func TestSyncMapRemove(t *testing.T) {
	c := newTestSyncMap()
	wg := new(sync.WaitGroup)
	for g := 0; g < testConcurrentGoroutines; g++ {
		g := g
		goroutine.WaitGroup(wg, func() {
			k := fmt.Sprintf("test_key_%d", g)
			v := fmt.Sprintf("test_value_%d", g)
			for i := 0; i < testConcurrentIterations; i++ {
				c.Set(k, v)
				c.Remove(k)
				_, found := c.Get(k)
				if found {
					t.Error("found")
					return
				}
			}
		})
	}
	wg.Wait()
}

func BenchmarkSyncMapSetGet(b *testing.B) {
	c := newTestSyncMap()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(testKey, testValue)
		c.Get(testKey)
	}
}

func BenchmarkSyncMapGet(b *testing.B) {
	c := newTestSyncMap()
	c.Set(testKey, testValue)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(testKey)
	}
}

func BenchmarkSyncMapParallelSetGet1(b *testing.B) {
	c := newTestSyncMap()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set(testKey, testValue)
			c.Get(testKey)
		}
	})
}

func BenchmarkSyncMapParallelSetGetN(b *testing.B) {
	c := newTestSyncMap()
	ct := newTestCounter()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		g := ct()
		var k interface{} = fmt.Sprintf("test_key_%d", g)
		var v interface{} = fmt.Sprintf("test_value_%d", g)
		for pb.Next() {
			c.Set(k, v)
			c.Get(k)
		}
	})
}

func BenchmarkSyncMapParallelGet(b *testing.B) {
	c := newTestSyncMap()
	c.Set(testKey, testValue)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get(testKey)
		}
	})
}

func newTestSyncMap() Cache {
	return New(
		Concurrent(),
	)
}
