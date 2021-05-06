package mongoutils

import (
	"context"

	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/mongo"
)

// ForceDisconnect disconnects the client with a context that is already canceled.
// It forces closes all opened connections.
func ForceDisconnect(clt *mongo.Client) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	return clt.Disconnect(ctx)
}

// NewForceDisconnect returns a closeutils.Err that calls ForceDisconnect.
func NewForceDisconnect(clt *mongo.Client) closeutils.Err {
	return func() error {
		return ForceDisconnect(clt)
	}
}
