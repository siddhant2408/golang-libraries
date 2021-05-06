// Package singleflight is a wrapper for golang.org/x/sync/singleflight.
//
// It provides tracing and better errors.
package singleflight

import (
	"context"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
	"golang.org/x/sync/singleflight" //nolint:depguard // The current package wraps singleflight.
)

// Group is a wrapper.
type Group struct {
	group
}

type group = singleflight.Group

// Do is a wrapper.
func (g *Group) Do(ctx context.Context, key string, fn func(context.Context) (interface{}, error)) (v interface{}, err error, shared bool) { //nolint:golint // The error is not the last returned item.
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "singleflight.do", &err)
	defer spanFinish()
	defer func() {
		// We need to add a new stack trace, because the "singleflight" algorithm is aggregating several calls together.
		// So if the error already contains a stack trace, it's not necessarily the one from the current call.
		err = errors.WithStack(err)
		span.SetTag("shared", shared)
	}()
	return g.group.Do(key, func() (interface{}, error) {
		return fn(ctx)
	})
}

// Result is a wrapper.
type Result = singleflight.Result
