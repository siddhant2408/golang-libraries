package ctxgrpcsrv

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/siddhant2408/golang-libraries/testutils"
	"google.golang.org/grpc"
)

func TestListenAndServe(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	port := getTestFreePort(t)
	addr := net.JoinHostPort("localhost", strconv.Itoa(port))
	srv := grpc.NewServer()
	err := ListenAndServe(ctx, addr, srv)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func getTestFreePort(tb testing.TB) int {
	tb.Helper()
	l := getTestFreeListener(tb)
	defer l.Close() //nolint:errcheck
	return l.Addr().(*net.TCPAddr).Port
}

func getTestFreeListener(tb testing.TB) *net.TCPListener {
	tb.Helper()
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		testutils.FatalErr(tb, err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		testutils.FatalErr(tb, err)
	}
	return l
}
