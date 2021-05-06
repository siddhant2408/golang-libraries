// Package bytespool provides a pool of []byte.
package bytespool

import "sync"

// Pool is a pool of []byte.
type Pool struct {
	p *sync.Pool
}

// New returns a Pool with length l.
func New(l int) *Pool {
	return &Pool{
		p: &sync.Pool{
			New: func() interface{} {
				c := containerPool.Get().(*container) //nolint:errcheck
				c.b = make([]byte, l)
				return c
			},
		},
	}
}

// Get gets a []byte from the Pool.
func (p *Pool) Get() []byte {
	c := p.p.Get().(*container) //nolint:errcheck
	b := c.b
	c.b = nil
	containerPool.Put(c)
	return b
}

// Put puts a []byte to the pool.
func (p *Pool) Put(b []byte) {
	c := containerPool.Get().(*container) //nolint:errcheck
	c.b = b
	p.p.Put(c)
}

// container allows to store []byte in a sync.Pool efficiently (as pointer type).
// It avoids to do memory allocation when we store a []byte in a sync.Pool.
type container struct {
	b []byte
}

var containerPool = &sync.Pool{
	New: func() interface{} {
		return new(container)
	},
}
