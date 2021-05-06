package goroutine

import (
	"sync"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestGo(t *testing.T) {
	var called testutils.CallCounter
	wait := Go(func() {
		called.Call()
	})
	wait()
	called.AssertCalled(t)
}

func TestWaitGroup(t *testing.T) {
	wg := new(sync.WaitGroup)
	var called testutils.CallCounter
	WaitGroup(wg, func() {
		called.Call()
	})
	wg.Wait()
	called.AssertCalled(t)
}

func TestRunN(t *testing.T) {
	var called testutils.CallCounter
	RunN(10, func() {
		called.Call()
	})
	called.AssertCount(t, 10)
}
