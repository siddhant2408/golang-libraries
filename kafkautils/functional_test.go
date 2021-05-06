package kafkautils

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/goroutine"
	"github.com/siddhant2408/golang-libraries/kafkatest"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestFunctionalProducerConsumer(t *testing.T) {
	ctx := context.Background()
	topic := kafkatest.Topic(t)
	cfg := &Config{
		Brokers: []string{kafkatest.GetBroker()},
	}
	w := cfg.NewWriter()
	w.BatchSize = 1
	defer w.Close() //nolint:errcheck
	p := &SimpleProducer{
		Writer: w.WriteMessages,
	}
	msgValue := []byte("test")
	msg := kafka.Message{
		Topic: topic,
		Value: msgValue,
	}
	err := p.Produce(ctx, msg)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	readerCfg := cfg.NewReaderConfig()
	readerCfg.GroupID = "test"
	readerCfg.Topic = topic
	readerCfg.MinBytes = 1
	readerCfg.MaxBytes = 1 << 20
	ctx, cancel := context.WithCancel(ctx)
	pr := func(ctx context.Context, msg kafka.Message) error {
		cancel()
		if !bytes.Equal(msg.Value, msgValue) {
			t.Error("not equal")
		}
		return nil
	}
	errFunc := func(ctx context.Context, err error) {
		testutils.ErrorErr(t, err)
	}
	RunConsumers(ctx, readerCfg, pr, 1, p.Produce, nil, errFunc)
}

func TestFunctionalProducerBatchConsumer(t *testing.T) {
	ctx := context.Background()
	topic := kafkatest.Topic(t)
	msgCount := 10
	cfg := &Config{
		Brokers: []string{kafkatest.GetBroker()},
	}
	w := cfg.NewWriter()
	w.BatchSize = 1
	defer w.Close() //nolint:errcheck
	p := &SimpleProducer{
		Writer: w.WriteMessages,
	}
	msgValue := []byte("test")
	var msgs []kafka.Message
	for i := 0; i < msgCount; i++ {
		msgs = append(msgs, kafka.Message{
			Topic: topic,
			Value: msgValue,
		})
	}
	err := p.Produce(ctx, msgs...)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	readerCfg := cfg.NewReaderConfig()
	readerCfg.GroupID = "test"
	readerCfg.Topic = topic
	readerCfg.MinBytes = 1
	readerCfg.MaxBytes = 1 << 20
	ctx, cancel := context.WithCancel(ctx)
	pr := func(ctx context.Context, msgs []kafka.Message) error {
		cancel()
		if len(msgs) != msgCount {
			t.Errorf("unexpected messages count: got %d, want %d", len(msgs), msgCount)
		}
		for _, msg := range msgs {
			if !bytes.Equal(msg.Value, msgValue) {
				t.Error("not equal")
			}
		}
		return nil
	}
	errFunc := func(ctx context.Context, err error) {
		testutils.ErrorErr(t, err)
	}
	RunBatchConsumers(ctx, readerCfg, pr, 1, msgCount, 0, errFunc)
}

func TestFunctionalWaitProducerConsumer(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	waitTopic := kafkatest.Topic(t)
	mainTopic := kafkatest.Topic(t)
	cfg := &Config{
		Brokers: []string{kafkatest.GetBroker()},
	}
	w := cfg.NewWriter()
	w.BatchSize = 1
	defer w.Close() //nolint:errcheck
	wsp := &SimpleProducer{
		Writer: w.WriteMessages,
	}
	wp := &WaitProducer{
		Producer: wsp.Produce,
		Wait:     1 * time.Millisecond,
	}
	msgValue := []byte("test")
	msg := kafka.Message{
		Topic: waitTopic,
		Value: msgValue,
	}
	err := wp.Produce(ctx, msg)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	wReaderCfg := cfg.NewReaderConfig()
	wReaderCfg.GroupID = "test"
	wReaderCfg.Topic = waitTopic
	wReaderCfg.MinBytes = 1
	wReaderCfg.MaxBytes = 1 << 20
	wReaderCfg.MaxWait = 1 * time.Second
	mp := &TopicProducer{
		Producer: (&SimpleProducer{
			Writer: w.WriteMessages,
		}).Produce,
		Topic: mainTopic,
	}
	waitWaitConsumers := goroutine.Go(func() {
		RunWaitConsumers(ctx, wReaderCfg, 1, mp.Produce, func(ctx context.Context, err error) {
			testutils.ErrorErr(t, err)
		})
	})
	defer waitWaitConsumers()
	mReaderCfg := cfg.NewReaderConfig()
	mReaderCfg.GroupID = "test"
	mReaderCfg.Topic = mainTopic
	mReaderCfg.MinBytes = 1
	mReaderCfg.MaxBytes = 1 << 20
	pr := func(ctx context.Context, msg kafka.Message) error {
		cancel()
		if !bytes.Equal(msg.Value, msgValue) {
			t.Error("not equal")
		}
		return nil
	}
	errFunc := func(ctx context.Context, err error) {
		testutils.ErrorErr(t, err)
	}
	RunConsumers(ctx, mReaderCfg, pr, 1, nil, nil, errFunc)
}
