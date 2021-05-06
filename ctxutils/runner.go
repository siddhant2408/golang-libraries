package ctxutils

import (
	"context"
	"time"
)

// Runner runs repetitively a function in a loop.
//
// If the context is cancelled, the runner is stopped.
// However the context cancellation is not propagated to the function (a new context is used).
type Runner struct {
	Func      Func
	ErrorFunc func(context.Context, error)
}

// Run runs the Runner.
func (r *Runner) Run(ctx context.Context) {
	r.run(ctx.Done())
}

func (r *Runner) run(done <-chan struct{}) {
	ctx := context.Background()
	ctx = setRunnerDoneContextValue(ctx, done)
	for !IsChannelClosed(done) {
		err := r.Func(ctx)
		if err != nil {
			r.ErrorFunc(ctx, err)
		}
	}
}

type runnerDoneContextKey struct{}

func setRunnerDoneContextValue(ctx context.Context, done <-chan struct{}) context.Context {
	return context.WithValue(ctx, runnerDoneContextKey{}, done)
}

func getRunnerDoneContextValue(ctx context.Context) <-chan struct{} {
	ch, ok := ctx.Value(runnerDoneContextKey{}).(<-chan struct{})
	if !ok {
		return nil
	}
	return ch
}

// NewRunnerFuncTicker returns a new Func that waits for a Ticker after each call.
//
// If an error is returned, it doesn't wait.
// If an error was returned at the previous call, it waits before the next call, without doing it.
//
// It is not concurrent safe.
func NewRunnerFuncTicker(f Func, tk *time.Ticker) Func {
	errored := false
	return func(ctx context.Context) error {
		if errored {
			runnerWait(ctx, tk.C)
			errored = false
			return nil
		}
		err := f(ctx)
		if err != nil {
			errored = true
			return err
		}
		runnerWait(ctx, tk.C)
		return nil
	}
}

// NewRunnerFuncWait returns a Func that waits for a specified duration after each call.
//
// If the returned duration is lower than or equal to zero, it doesn't wait.
// If an error is returned, it doesn't wait.
func NewRunnerFuncWait(f func(context.Context) (time.Duration, error)) Func {
	return func(ctx context.Context) error {
		wait, err := f(ctx)
		if err != nil {
			return err
		}
		if wait <= 0 {
			return nil
		}
		tm := time.NewTimer(wait)
		defer tm.Stop()
		runnerWait(ctx, tm.C)
		return nil
	}
}

func runnerWait(ctx context.Context, tmc <-chan time.Time) {
	done := getRunnerDoneContextValue(ctx)
	select {
	case <-done:
	case <-tmc:
	}
}
