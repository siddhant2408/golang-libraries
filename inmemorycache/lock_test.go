package inmemorycache

import (
	"fmt"
	"sync"
	"testing"

	"github.com/siddhant2408/golang-libraries/goroutine"
)

const (
	testConcurrentGoroutines = 64
	testConcurrentIterations = 1000
)

func TestLockSetGet(t *testing.T) {
	c := newTestLock()
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

func TestLockGet(t *testing.T) {
	c := newTestLock()
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

func TestLockRemove(t *testing.T) {
	c := newTestLock()
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

func BenchmarkLockSetGet(b *testing.B) {
	c := newTestLock()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(testKey, testValue)
		c.Get(testKey)
	}
}

func BenchmarkLockGet(b *testing.B) {
	c := newTestLock()
	c.Set(testKey, testValue)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(testKey)
	}
}

func BenchmarkLockParallelSetGet1(b *testing.B) {
	c := newTestLock()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set(testKey, testValue)
			c.Get(testKey)
		}
	})
}

func BenchmarkLockParallelSetGetN(b *testing.B) {
	c := newTestLock()
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

func BenchmarkLockParallelGet(b *testing.B) {
	c := newTestLock()
	c.Set(testKey, testValue)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get(testKey)
		}
	})
}

func newTestLock() Cache {
	return New(
		disableSyncMap(), // Must disable sync map, otherwise the lock cache is not used.
		Concurrent(),
	)
}

func newTestCounter() func() int {
	var mu sync.Mutex
	var i int
	return func() int {
		mu.Lock()
		v := i
		i++
		mu.Unlock()
		return v
	}
}
