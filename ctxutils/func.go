package ctxutils

import (
	"context"
	"sync"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/goroutine"
)

// Func represents a function to run.
type Func func(context.Context) error

// Funcs represents a group of functions to run.
// Each function is identified by a unique name (map key).
type Funcs map[string]Func

// RunFuncs runs a group of functions.
// Each function runs in a goroutine.
// It returns once all functions returns.
// If a function returns an error, the context passed to all functions is canceled.
// Only the first error is returned.
func RunFuncs(ctx context.Context, fs Funcs) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errCh := make(chan error, 1)
	wg := new(sync.WaitGroup)
	for name, f := range fs {
		name, f := name, f
		goroutine.WaitGroup(wg, func() {
			err := f(ctx)
			if err == nil {
				return
			}
			cancel()
			err = errors.Wrapf(err, "run %q", name)
			select {
			case errCh <- err:
			default:
			}
		})
	}
	wg.Wait()
	select {
	case err := <-errCh:
		return errors.WithStack(err)
	default:
		return nil
	}
}
