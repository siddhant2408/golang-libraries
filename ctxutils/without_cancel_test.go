package ctxutils

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestWithoutCancel(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	ctx = WithoutCancel(ctx)
	done := IsDone(ctx)
	if done {
		t.Fatal("done")
	}
	_, ok := ctx.Deadline()
	if ok {
		t.Fatal("deadline")
	}
	err := ctx.Err()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
