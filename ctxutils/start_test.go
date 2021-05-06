package ctxutils

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestStart(t *testing.T) {
	ctx := context.Background()
	var called testutils.CallCounter
	cl := Start(ctx, func(ctx context.Context) {
		called.Call()
		<-ctx.Done()
	})
	cl()
	called.AssertCalled(t)
}
