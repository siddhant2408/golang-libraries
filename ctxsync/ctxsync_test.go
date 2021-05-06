package ctxsync

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestMutexLock(t *testing.T) {
	m := &Mutex{}
	for i := 0; i < 10; i++ {
		m.Lock()
		m.Unlock() //nolint:gocritic,staticcheck // This test does nothing.
	}
}

func TestMutexLockCtx(t *testing.T) {
	ctx := context.Background()
	m := &Mutex{}
	for i := 0; i < 10; i++ {
		err := m.LockCtx(ctx)
		if err != nil {
			testutils.FatalErr(t, err)
		}
		m.Unlock()
	}
}

func TestMutexLockCtxErrorContextDone1(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	m := &Mutex{}
	err := m.LockCtx(ctx)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestMutexLockCtxErrorContextDone2(t *testing.T) {
	ctx := context.Background()
	m := &Mutex{}
	err := m.LockCtx(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()
	err = m.LockCtx(ctx)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestMutexUnlockPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("no panic")
		}
	}()
	m := &Mutex{}
	m.Unlock()
}

func BenchmarkMutexLock(b *testing.B) {
	m := &Mutex{}
	for i := 0; i < b.N; i++ {
		m.Lock()
		m.Unlock() //nolint:gocritic,staticcheck // This test does nothing.
	}
}

func BenchmarkMutexLockCtx(b *testing.B) {
	ctx := context.Background()
	m := &Mutex{}
	for i := 0; i < b.N; i++ {
		err := m.LockCtx(ctx)
		if err != nil {
			testutils.FatalErr(b, err)
		}
		m.Unlock()
	}
}

func BenchmarkStdlibMutex(b *testing.B) {
	m := &sync.Mutex{}
	for i := 0; i < b.N; i++ {
		m.Lock()
		m.Unlock() //nolint:gocritic,staticcheck // This test does nothing.
	}
}
