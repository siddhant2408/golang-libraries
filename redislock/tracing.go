package redislock

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

func startTraceSpan(pctx *context.Context, operationName string, perr *error) (opentracing.Span, closeutils.F) {
	return tracingutils.StartChildSpan(pctx, "redislock."+operationName, perr)
}

func setTraceSpanTag(span opentracing.Span, key string, val interface{}) {
	span.SetTag("redislock."+key, val)
}
