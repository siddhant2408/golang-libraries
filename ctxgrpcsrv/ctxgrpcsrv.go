// Package ctxgrpcsrv provides a helper to stop agRPC Server if a Context is canceled.
// It returns no error if the Context is canceled.
package ctxgrpcsrv

import (
	"context"
	"net"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/goroutine"
	"google.golang.org/grpc"
)

// ListenAndServe listens on addr and calls Serve.
func ListenAndServe(ctx context.Context, addr string, srv *grpc.Server) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "listen")
	}
	return Serve(ctx, l, srv)
}

// Serve is a replacement for grpc.Server.Serve.
func Serve(ctx context.Context, l net.Listener, srv *grpc.Server) error {
	errCh := make(chan error)
	wait := goroutine.Go(func() {
		err := srv.Serve(l)
		select {
		case errCh <- err:
		case <-ctx.Done():
		}
	})
	defer wait()
	select {
	case err := <-errCh:
		return errors.WithStack(err)
	case <-ctx.Done():
		srv.GracefulStop()
		return nil
	}
}
