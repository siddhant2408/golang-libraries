package ctxutils

import (
	"context"
	"testing"
)

func TestIsDoneTrue(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	done := IsDone(ctx)
	if !done {
		t.Fatalf("unexpected done: got %t, want %t", done, true)
	}
}

func TestIsDoneFalse(t *testing.T) {
	ctx := context.Background()
	done := IsDone(ctx)
	if done {
		t.Fatalf("unexpected done: got %t, want %t", done, false)
	}
}
