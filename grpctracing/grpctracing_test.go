package grpctracing

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
	"google.golang.org/grpc"
)

func TestUnaryServerInterceptor(t *testing.T) {
	ctx := context.Background()
	req := map[string]interface{}{
		"foo": "bar",
	}
	info := &grpc.UnaryServerInfo{
		FullMethod: "test",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return map[string]interface{}{
			"foo": "bar",
		}, nil
	}
	resp, err := UnaryServerInterceptor(ctx, req, info, handler)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if resp == nil {
		t.Fatal("nil")
	}
}
