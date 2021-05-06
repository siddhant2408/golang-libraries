package ctxutils

import (
	"context"
	"strconv"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestRunFuncs(t *testing.T) {
	for i := 0; i <= 10; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			fs := make(Funcs)
			for j := 1; j <= i; j++ {
				fs[strconv.Itoa(j)] = func(ctx context.Context) error {
					return nil
				}
			}
			err := RunFuncs(context.Background(), fs)
			if err != nil {
				testutils.FatalErr(t, err)
			}
		})
	}
}

func TestRunFuncsErrorSingle(t *testing.T) {
	for i := 1; i <= 10; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			fs := make(Funcs)
			for j := 1; j <= i; j++ {
				var f Func
				if j == 1 {
					f = func(ctx context.Context) error {
						return errors.New("error")
					}
				} else {
					f = func(ctx context.Context) error {
						return nil
					}
				}
				fs[strconv.Itoa(j)] = f
			}
			err := RunFuncs(context.Background(), fs)
			if err == nil {
				t.Fatal("no error")
			}
		})
	}
}

func TestRunFuncsErrorAll(t *testing.T) {
	for i := 1; i <= 10; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			fs := make(Funcs)
			for j := 1; j <= i; j++ {
				fs[strconv.Itoa(j)] = func(ctx context.Context) error {
					return errors.New("error")
				}
			}
			err := RunFuncs(context.Background(), fs)
			if err == nil {
				t.Fatal("no error")
			}
		})
	}
}

func TestRunFuncsContextCanceled(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	err := RunFuncs(ctx, Funcs{
		"test": func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		},
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
