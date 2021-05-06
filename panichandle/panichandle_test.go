package panichandle

import (
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestDefaultHandler(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("no panic")
		}
	}()
	DefaultHandler("test")
}

func TestRecover(t *testing.T) {
	Recover()
}

func TestRecoverPanic(t *testing.T) {
	defer restoreDefaultHandler()
	var called testutils.CallCounter
	Handler = func(r interface{}) {
		called.Call()
		if r == nil {
			t.Fatal("nil")
		}
	}
	defer func() {
		called.AssertCalled(t)
	}()
	defer Recover()
	panic("test")
}

func restoreDefaultHandler() {
	Handler = DefaultHandler
}
