package ctxutils

import (
	"context"
)

// IsDone returns true if the context is done, false otherwise.
func IsDone(ctx context.Context) bool {
	return IsChannelClosed(ctx.Done())
}

// IsChannelClosed returns true if the channel is closed, false otherwise.
func IsChannelClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
