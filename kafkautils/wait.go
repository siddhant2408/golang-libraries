package kafkautils

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/ctxutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/timeutils"
)

const (
	// WaitUntilHeader is the header containing the "wait until" date (RFC3339).
	WaitUntilHeader       = "wait-until"
	waitUntilHeaderLayout = time.RFC3339Nano
)

// WaitProducer wraps a Producer and adds a "wait-until" header, relative to the current date.
type WaitProducer struct {
	Producer
	Wait time.Duration
}

// Produce adds a "wait-until" header and produces the messages.
func (p *WaitProducer) Produce(ctx context.Context, msgs ...kafka.Message) error {
	waitUntil := timeutils.Now().Add(p.Wait)
	waitUntilStr := waitUntil.Format(waitUntilHeaderLayout)
	waitUntilBytes := []byte(waitUntilStr)
	for i, msg := range msgs {
		msg.Headers = SetHeader(msg.Headers, WaitUntilHeader, waitUntilBytes)
		msgs[i] = msg
	}
	return p.Producer(ctx, msgs...)
}

// WaitConsumer allows to wait before processing a message.
//
// The steps are:
//  - fetch a message
//  - wait until the date defined by the "wait-until" header
//  - produce the message to Producer ("main" topic)
//  - commit the message
//
// For each message, it waits until the date specified in the "wait-until" header.
// If the header is not defined or is invalid, it doesn't wait.
//
// Warning: each consumer should read a single partition only.
// Otherwise the wait delay might be wrongly interpreted.
// The topic should have a single partition.
type WaitConsumer struct {
	Producer Producer
}

// Consume consumes messages from a reader.
func (c *WaitConsumer) Consume(ctx context.Context, r FetchCommitter) error {
	for !ctxutils.IsDone(ctx) {
		err := c.consume(ctx, r)
		if err != nil {
			return errors.Wrap(err, "wait consumer")
		}
	}
	return nil
}

func (c *WaitConsumer) consume(ctx context.Context, r FetchCommitter) error {
	msg, err := r.FetchMessage(ctx)
	if err != nil {
		if ctxutils.IsDone(ctx) {
			return nil
		}
		return errors.Wrap(err, "fetch")
	}
	c.wait(ctx, msg)
	if ctxutils.IsDone(ctx) {
		return nil
	}
	ctx = context.Background() // Don't want to be interrupted
	newMsg := CopyMessage(msg)
	newMsg.Headers = DeleteHeader(newMsg.Headers, WaitUntilHeader)
	err = c.Producer(ctx, newMsg)
	if err != nil {
		return errors.Wrap(err, "produce")
	}
	err = r.CommitMessages(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "commit")
	}
	return nil
}

func (c *WaitConsumer) wait(ctx context.Context, msg kafka.Message) {
	waitUntilBytes, ok := GetHeader(msg.Headers, WaitUntilHeader)
	if !ok {
		// If the header is not defined: don't wait.
		return
	}
	waitUntilStr := string(waitUntilBytes)
	waitUntil, err := time.Parse(waitUntilHeaderLayout, waitUntilStr)
	if err != nil {
		// If the date format is invalid: don't wait.
		return
	}
	dur := timeutils.Until(waitUntil)
	if dur <= 0 {
		// If the date is in the past: don't wait.
		return
	}
	tm := time.NewTimer(dur)
	defer tm.Stop()
	// Wait until the date or the context is done.
	select {
	case <-tm.C:
	case <-ctx.Done():
	}
}

// RunWaitConsumers runs wait consumers.
func RunWaitConsumers(ctx context.Context, readerCfg kafka.ReaderConfig, count int, producer Producer, errFunc func(context.Context, error)) {
	c := &WaitConsumer{
		Producer: producer,
	}
	ConsumeReader(ctx, readerCfg, c.Consume, errFunc)
}
