package ctxutils

import (
	"context"

	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/goroutine"
)

// Start executes the given function in a new goroutine.
// The close function cancels the context passed to the function, and wait until the goroutine exits.
func Start(ctx context.Context, f func(context.Context)) closeutils.F {
	ctx, cancel := context.WithCancel(ctx)
	wait := goroutine.Go(func() {
		f(ctx)
	})
	cl := func() {
		cancel()
		wait()
	}
	return cl
}
