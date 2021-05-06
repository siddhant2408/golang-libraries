package testutils

import (
	"sync/atomic"
	"testing"
)

// CallCounter counts the number of call.
type CallCounter struct {
	calls int64
}

// Call increment the counter by 1.
func (c *CallCounter) Call() {
	atomic.AddInt64(&c.calls, 1)
}

// AssertCalled is an alias for AssertCalledFatal.
func (c *CallCounter) AssertCalled(tb testing.TB) {
	tb.Helper()
	c.AssertCalledFatal(tb)
}

// AssertCalledFatal checks that the counter has been called at least once, and calls fatal if it's not the case.
func (c *CallCounter) AssertCalledFatal(tb testing.TB) {
	tb.Helper()
	c.assertCalled(tb, tb.Fatal)
}

// AssertCalledError checks that the counter has been called at least once, and calls error if it's not the case.
func (c *CallCounter) AssertCalledError(tb testing.TB) {
	tb.Helper()
	c.assertCalled(tb, tb.Error)
}

// AssertCalledSkip checks that the counter has been called at least once, and calls skip if it's not the case.
func (c *CallCounter) AssertCalledSkip(tb testing.TB) {
	tb.Helper()
	c.assertCalled(tb, tb.Skip)
}

func (c *CallCounter) assertCalled(tb testing.TB, f func(args ...interface{})) {
	tb.Helper()
	if c.calls < 1 {
		f("not called")
	}
}

// AssertNotCalled is an alias for AssertNotCalledFatal.
func (c *CallCounter) AssertNotCalled(tb testing.TB) {
	tb.Helper()
	c.AssertNotCalledFatal(tb)
}

// AssertNotCalledFatal checks that the counter has not been called, and calls fatal if it's not the case.
func (c *CallCounter) AssertNotCalledFatal(tb testing.TB) {
	tb.Helper()
	c.assertNotCalled(tb, tb.Fatal)
}

// AssertNotCalledError checks that the counter has not been called, and calls error if it's not the case.
func (c *CallCounter) AssertNotCalledError(tb testing.TB) {
	tb.Helper()
	c.assertNotCalled(tb, tb.Error)
}

// AssertNotCalledSkip checks that the counter has not been called, and calls skip if it's not the case.
func (c *CallCounter) AssertNotCalledSkip(tb testing.TB) {
	tb.Helper()
	c.assertNotCalled(tb, tb.Skip)
}

func (c *CallCounter) assertNotCalled(tb testing.TB, f func(args ...interface{})) {
	tb.Helper()
	if c.calls > 0 {
		f("called")
	}
}

// AssertCount is an alias for AssertCountFatal.
func (c *CallCounter) AssertCount(tb testing.TB, expected int64) {
	tb.Helper()
	c.AssertCountFatal(tb, expected)
}

// AssertCountFatal checks that the call count is equal to the expected value, and calls fatal if it's not the case.
func (c *CallCounter) AssertCountFatal(tb testing.TB, expected int64) {
	tb.Helper()
	c.assertCount(tb, tb.Fatalf, expected)
}

// AssertCountError checks that the call count is equal to the expected value, and calls error if it's not the case.
func (c *CallCounter) AssertCountError(tb testing.TB, expected int64) {
	tb.Helper()
	c.assertCount(tb, tb.Errorf, expected)
}

// AssertCountSkip checks that the call count is equal to the expected value, and calls skip if it's not the case.
func (c *CallCounter) AssertCountSkip(tb testing.TB, expected int64) {
	tb.Helper()
	c.assertCount(tb, tb.Skipf, expected)
}

func (c *CallCounter) assertCount(tb testing.TB, f func(format string, args ...interface{}), expected int64) {
	tb.Helper()
	if c.calls != expected {
		f("unexpected call count: got %d, want %d", c.calls, expected)
	}
}
