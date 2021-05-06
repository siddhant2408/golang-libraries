package kafkautils

import (
	"context"
	"sort"

	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

// Producer produces messages.
type Producer func(context.Context, ...kafka.Message) error

// SimpleProducer produce messages.
type SimpleProducer struct {
	Writer func(context.Context, ...kafka.Message) error
}

// Produce produces messages.
func (p *SimpleProducer) Produce(ctx context.Context, msgs ...kafka.Message) (err error) {
	span, spanFinish := startTraceChildSpan(&ctx, "producer", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, tracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.SpanTypeMessageProducer)
	opentracing_ext.SpanKindProducer.Set(span)
	setTraceSpanTag(span, "messages.count", len(msgs))
	if len(msgs) == 1 {
		setTraceSpanTagsMessage(span, msgs[0])
	}
	err = p.Writer(ctx, msgs...)
	if err != nil {
		err = p.wrapWriteErrors(err)
		if len(msgs) == 1 {
			err = wrapErrorValueMessage(err, msgs[0])
		}
		return errors.Wrap(err, "writer")
	}
	return nil
}

func (p *SimpleProducer) wrapWriteErrors(err error) error {
	var werrs kafka.WriteErrors
	ok := errors.As(err, &werrs)
	if !ok {
		return err
	}
	errMsgSet := make(map[string]struct{})
	for _, werr := range werrs {
		if werr == nil {
			continue
		}
		errMsg := werr.Error()
		errMsgSet[errMsg] = struct{}{}
	}
	errMsgs := make([]string, 0, len(errMsgSet))
	for errMsg := range errMsgSet {
		errMsgs = append(errMsgs, errMsg)
	}
	sort.Strings(errMsgs) // For stable output.
	err = wrapErrorValue(err, "write_errors", errMsgs)
	return err
}

// TopicProducer is a Producer that overwrites the `kafka.Message.Topic` field.
type TopicProducer struct {
	Producer
	Topic string
}

// Produce produces messages.
func (p *TopicProducer) Produce(ctx context.Context, msgs ...kafka.Message) (err error) {
	for i, msg := range msgs {
		msg.Topic = p.Topic
		msgs[i] = msg
	}
	return p.Producer(ctx, msgs...)
}
