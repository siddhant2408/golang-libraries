package bufpool

import (
	"bytes"
	"testing"
)

func Test(t *testing.T) {
	buf := Get()
	Put(buf)
}

const testBenchData = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.\n"

var benchmarkResult string

func BenchmarkBufPool(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		buf := Get()
		for j := 0; j < 10; j++ {
			buf.WriteString(testBenchData)
		}
		s = buf.String()
		Put(buf)
	}
	benchmarkResult = s
}

func BenchmarkBufWithoutPool(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		for j := 0; j < 10; j++ {
			buf.WriteString(testBenchData)
		}
		s = buf.String()
	}
	benchmarkResult = s
}
