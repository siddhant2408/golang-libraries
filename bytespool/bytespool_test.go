package bytespool

import (
	"testing"
)

const testLength = 32 << 10

func Test(t *testing.T) {
	p := New(testLength)
	for i := 0; i < 10; i++ {
		buf := p.Get()
		if len(buf) != testLength {
			t.Fatalf("unexpected length: got %d, want %d", len(buf), testLength)
		}
		p.Put(buf)
	}
}

func TestAllocs(t *testing.T) {
	p := New(testLength)
	allocs := testing.AllocsPerRun(1000, func() {
		buf := p.Get()
		p.Put(buf)
	})
	if allocs != 0 {
		t.Fatalf("unexpected allocs: got %f, want %d", allocs, 0)
	}
}

func Benchmark(b *testing.B) {
	p := New(testLength)
	for i := 0; i < b.N; i++ {
		buf := p.Get()
		p.Put(buf)
	}
}
