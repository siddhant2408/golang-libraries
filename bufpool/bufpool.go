// Package bufpool provides a sync.Pool of bytes.Buffer.
package bufpool

import (
	"bytes"
	"sync"
)

var defaultPool = New()

// Get calls Get on the default Pool.
func Get() *bytes.Buffer {
	return defaultPool.Get()
}

// Put calls Put on the default Pool.
func Put(buf *bytes.Buffer) {
	defaultPool.Put(buf)
}

// Pool is a pool of *bytes.Buffer.
type Pool struct {
	p *sync.Pool

	// MaxCap defines the maximum capacity accepted for recycled buffer.
	// If Put() is called with a buffer larger than this value, it's discarded.
	// See https://github.com/golang/go/issues/23199 .
	// 0 means there is no maximum capacity.
	MaxCap int
}

const maxCapDefault = 1 << 16 // 64 KiB

// New returns a new Pool.
func New() *Pool {
	return &Pool{
		p: &sync.Pool{
			New: syncPoolNew,
		},
		MaxCap: maxCapDefault,
	}
}

func syncPoolNew() interface{} {
	return new(bytes.Buffer)
}

// Get gets a buffer from the Pool, resets it and returns it.
func (p *Pool) Get() *bytes.Buffer {
	buf := p.p.Get().(*bytes.Buffer) //nolint:errcheck
	// we rest the buffer because we don't want to reuse the underlying / previously stored values
	buf.Reset()
	return buf
}

// Put puts the buffer to the Pool.
// WARNING: the call MUST NOT reuse the buffer's content after this call.
func (p *Pool) Put(buf *bytes.Buffer) {
	if p.MaxCap != 0 && buf.Cap() <= p.MaxCap {
		p.p.Put(buf)
	}
}
