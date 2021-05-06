package kafkautils

import (
	"context"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/siddhant2408/golang-libraries/timeutils"
)

func TestSimpleProducer(t *testing.T) {
	ctx := context.Background()
	var wCalled testutils.CallCounter
	w := func(context.Context, ...kafka.Message) error {
		wCalled.Call()
		return nil
	}
	p := &SimpleProducer{
		Writer: w,
	}
	msg := kafka.Message{}
	err := p.Produce(ctx, msg)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	wCalled.AssertCalled(t)
}

func TestSimpleProducerError(t *testing.T) {
	ctx := context.Background()
	w := func(context.Context, ...kafka.Message) error {
		return errors.New("error")
	}
	p := &SimpleProducer{
		Writer: w,
	}
	msg := kafka.Message{
		Key:   []byte("test"),
		Value: []byte("test"),
		Time:  timeutils.Now(),
	}
	err := p.Produce(ctx, msg)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSimpleProducerErrorWrite(t *testing.T) {
	ctx := context.Background()
	w := func(context.Context, ...kafka.Message) error {
		return kafka.WriteErrors{
			kafka.UnknownTopicOrPartition,
			kafka.InvalidMessage,
			nil,
			kafka.InvalidMessage,
		}
	}
	p := &SimpleProducer{
		Writer: w,
	}
	msg := kafka.Message{}
	err := p.Produce(ctx, msg)
	if err == nil {
		t.Fatal("no error")
	}
	errMsgs, ok := errors.Values(err)["kafka.write_errors"].([]string)
	if !ok {
		t.Fatal("no write error messages")
	}
	expectedErrMsgs := []string{
		kafka.InvalidMessage.Error(),
		kafka.UnknownTopicOrPartition.Error(),
	}
	testutils.Compare(t, "unexpected write error messages", errMsgs, expectedErrMsgs)
}

func TestTopicProducer(t *testing.T) {
	ctx := context.Background()
	topic := "test"
	var pCalled testutils.CallCounter
	p := func(ctx context.Context, msgs ...kafka.Message) (err error) {
		pCalled.Call()
		for _, msg := range msgs {
			if msg.Topic != topic {
				t.Fatalf("unexpected topic: got %q, want %q", msg.Topic, topic)
			}
		}
		return nil
	}
	tp := &TopicProducer{
		Producer: p,
		Topic:    topic,
	}
	msg := kafka.Message{
		Topic: topic,
	}
	err := tp.Produce(ctx, msg)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	pCalled.AssertCalled(t)
}

func TestTopicProducerError(t *testing.T) {
	ctx := context.Background()
	topic := "test"
	p := func(ctx context.Context, msgs ...kafka.Message) (err error) {
		return errors.New("error")
	}
	tp := &TopicProducer{
		Producer: p,
		Topic:    topic,
	}
	msg := kafka.Message{
		Topic: topic,
	}
	err := tp.Produce(ctx, msg)
	if err == nil {
		t.Fatal("no error")
	}
}
