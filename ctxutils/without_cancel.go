package ctxutils

import (
	"context"
	"time"
)

// WithoutCancel returns a new child context that is not cancelled when the parent is cancelled.
func WithoutCancel(ctx context.Context) context.Context {
	return &withoutCancel{
		Context: ctx,
	}
}

type withoutCancel struct {
	context.Context
}

func (ctx *withoutCancel) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (ctx *withoutCancel) Done() <-chan struct{} {
	return nil
}

func (ctx *withoutCancel) Err() error {
	return nil
}
