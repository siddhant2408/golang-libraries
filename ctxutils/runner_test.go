package ctxutils

import (
	"context"
	"testing"
	"time"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestRunner(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	r := &Runner{
		Func: func(context.Context) error {
			cancel()
			return nil
		},
	}
	r.Run(ctx)
}

func TestRunnerError(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	r := &Runner{
		Func: func(context.Context) error {
			return errors.New("error")
		},
		ErrorFunc: func(context.Context, error) {
			cancel()
		},
	}
	r.Run(ctx)
}

func TestNewRunnerFuncTicker(t *testing.T) {
	ctx := context.Background()
	var called testutils.CallCounter
	tk := time.NewTicker(1 * time.Microsecond)
	defer tk.Stop()
	f := NewRunnerFuncTicker(func(context.Context) error {
		called.Call()
		return nil
	}, tk)
	err := f(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	called.AssertCalled(t)
}

func TestNewRunnerFuncTickerDone(t *testing.T) {
	ctx := context.Background()
	done := make(chan struct{})
	close(done)
	ctx = setRunnerDoneContextValue(ctx, done)
	var called testutils.CallCounter
	tk := time.NewTicker(1 * time.Hour)
	defer tk.Stop()
	f := NewRunnerFuncTicker(func(context.Context) error {
		called.Call()
		return nil
	}, tk)
	err := f(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	called.AssertCalled(t)
}

func TestNewRunnerFuncTickerError(t *testing.T) {
	ctx := context.Background()
	var called testutils.CallCounter
	tk := time.NewTicker(1 * time.Microsecond)
	defer tk.Stop()
	f := NewRunnerFuncTicker(func(context.Context) error {
		called.Call()
		return errors.New("error")
	}, tk)
	err := f(ctx)
	if err == nil {
		t.Fatal("no error")
	}
	called.AssertCount(t, 1)
	err = f(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	called.AssertCount(t, 1)
	err = f(ctx)
	if err == nil {
		t.Fatal("no error")
	}
	called.AssertCount(t, 2)
}

func TestNewRunnerFuncWait(t *testing.T) {
	ctx := context.Background()
	var called testutils.CallCounter
	f := NewRunnerFuncWait(func(context.Context) (time.Duration, error) {
		called.Call()
		return 1 * time.Microsecond, nil
	})
	err := f(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	called.AssertCalled(t)
}

func TestNewRunnerFuncWaitZero(t *testing.T) {
	ctx := context.Background()
	var called testutils.CallCounter
	f := NewRunnerFuncWait(func(context.Context) (time.Duration, error) {
		called.Call()
		return 0, nil
	})
	err := f(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	called.AssertCalled(t)
}

func TestNewRunnerFuncWaitDone(t *testing.T) {
	ctx := context.Background()
	done := make(chan struct{})
	close(done)
	ctx = setRunnerDoneContextValue(ctx, done)
	var called testutils.CallCounter
	f := NewRunnerFuncWait(func(context.Context) (time.Duration, error) {
		called.Call()
		return 1 * time.Hour, nil
	})
	err := f(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	called.AssertCalled(t)
}

func TestNewRunnerFuncWaitError(t *testing.T) {
	ctx := context.Background()
	var called testutils.CallCounter
	f := NewRunnerFuncWait(func(context.Context) (time.Duration, error) {
		called.Call()
		return 0, errors.New("error")
	})
	err := f(ctx)
	if err == nil {
		t.Fatal("no error")
	}
	called.AssertCalled(t)
}
