// Package grpctracing provides gRPC tracing related utilities.
package grpctracing

import (
	"context"

	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/siddhant2408/golang-libraries/tracingutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// UnaryServerInterceptor is a grpc.UnaryServerInterceptor.
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	span, spanFinish := tracingutils.StartRootSpan(&ctx, "grpc.server", &err)
	defer spanFinish()
	tracingutils.SetSpanType(span, tracingutils.AppTypeRPC)
	opentracing_ext.SpanKindRPCServer.Set(span)
	tracingutils.SetSpanResourceName(span, info.FullMethod)
	span.SetTag("grpc.method.name", info.FullMethod)
	span.SetTag("grpc.method.kind", "unary")
	pe, ok := peer.FromContext(ctx)
	if ok {
		span.SetTag("grpc.peer.addr", pe.Addr.String())
	}
	pm, ok := req.(proto.Message)
	if ok {
		var b []byte
		b, err = protojson.Marshal(pm)
		if err == nil {
			span.SetTag("grpc.request", string(b))
		}
	}
	md, _ := metadata.FromIncomingContext(ctx)
	for k, v := range md {
		span.SetTag("grpc.metadata."+k, v)
	}
	resp, err = handler(ctx, req)
	code := status.Code(err)
	span.SetTag("grpc.code", code.String())
	return resp, err
}
