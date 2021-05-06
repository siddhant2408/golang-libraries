package kafkautils

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

const (
	tracingExternalServiceName = "go-kafka"
)

func startTraceRootSpan(pctx *context.Context, op string, perr *error) (opentracing.Span, closeutils.F) {
	return tracingutils.StartRootSpan(pctx, "kafka."+op, perr)
}

func startTraceChildSpan(pctx *context.Context, op string, perr *error) (opentracing.Span, closeutils.F) {
	return tracingutils.StartChildSpan(pctx, "kafka."+op, perr)
}

func setTraceSpanTag(span opentracing.Span, key string, val interface{}) {
	span.SetTag("kafka."+key, val)
}

func setTraceSpanTagsMessage(span opentracing.Span, msg kafka.Message) {
	if len(msg.Key) > 0 {
		setTraceSpanTag(span, "message.key", bytesTruncateConvert(msg.Key))
	}
	if len(msg.Value) > 0 {
		setTraceSpanTag(span, "message.value", bytesTruncateConvert(msg.Value))
	}
	if !msg.Time.IsZero() {
		setTraceSpanTag(span, "message.time", msg.Time.Format(time.RFC3339Nano))
	}
	// TODO headers
}

func setTraceSpanTagsMessageConsumer(span opentracing.Span, msg kafka.Message) {
	setTraceSpanTagsMessage(span, msg)
	setTraceSpanTag(span, "message.topic", msg.Topic)
	setTraceSpanTag(span, "message.partition", msg.Partition)
	setTraceSpanTag(span, "message.offset", msg.Offset)
}
