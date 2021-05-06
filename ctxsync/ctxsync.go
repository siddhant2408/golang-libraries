// Package ctxsync provides sync with context support.
package ctxsync

import (
	"context"
	"sync"
)

// Mutex is similar to a sync.Mutex and supports context cancellation.
type Mutex struct {
	once sync.Once
	ch   chan struct{}
}

// Lock locks the mutex.
func (m *Mutex) Lock() {
	m.init()
	m.ch <- struct{}{}
}

// LockCtx locks the mutex with a context.
//
// If the context is done before the lock is acquired, the context's error is returned.
func (m *Mutex) LockCtx(ctx context.Context) error {
	m.init()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	select {
	case m.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Unlock unlock the mutex.
func (m *Mutex) Unlock() {
	m.init()
	select {
	case <-m.ch:
	default:
		panic("unlock of unlocked mutex")
	}
}

func (m *Mutex) init() {
	m.once.Do(func() {
		m.ch = make(chan struct{}, 1)
	})
}
