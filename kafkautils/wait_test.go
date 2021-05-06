package kafkautils

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/siddhant2408/golang-libraries/timeutils"
)

func TestWaitProducer(t *testing.T) {
	ctx := context.Background()
	wait := 1 * time.Minute
	msgs := []kafka.Message{
		{
			Value: []byte("value 1"),
		},
		{
			Value: []byte("value 2"),
		},
	}
	expectedMsgs := []kafka.Message{
		{
			Value: []byte("value 1"),
			Headers: []kafka.Header{
				{
					Key:   WaitUntilHeader,
					Value: []byte(timeutils.Now().Add(wait).Format(waitUntilHeaderLayout)),
				},
			},
		},
		{
			Value: []byte("value 2"),
			Headers: []kafka.Header{
				{
					Key:   WaitUntilHeader,
					Value: []byte(timeutils.Now().Add(wait).Format(waitUntilHeaderLayout)),
				},
			},
		},
	}
	var pCalled testutils.CallCounter
	p := func(ctx context.Context, msgs ...kafka.Message) error {
		pCalled.Call()
		testutils.Compare(t, "unexpected messages", msgs, expectedMsgs)
		return nil
	}
	wp := &WaitProducer{
		Producer: p,
		Wait:     wait,
	}
	err := wp.Produce(ctx, msgs...)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	pCalled.AssertCalled(t)
}

func TestWaitProducerError(t *testing.T) {
	ctx := context.Background()
	wait := 1 * time.Minute
	msgs := []kafka.Message{
		{
			Value: []byte("value 1"),
		},
		{
			Value: []byte("value 2"),
		},
	}
	p := func(ctx context.Context, msgs ...kafka.Message) error {
		return errors.New("error")
	}
	wp := &WaitProducer{
		Producer: p,
		Wait:     wait,
	}
	err := wp.Produce(ctx, msgs...)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestWaitConsumer(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	msg := kafka.Message{
		Value: []byte("test"),
		Headers: []kafka.Header{
			{
				Key:   WaitUntilHeader,
				Value: []byte(timeutils.Now().Add(1 * time.Microsecond).Format(waitUntilHeaderLayout)),
			},
		},
	}
	expectedMsgs := []kafka.Message{
		{
			Value: []byte("test"),
		},
	}
	p := func(ctx context.Context, msgs ...kafka.Message) error {
		testutils.Compare(t, "unexpected message", msgs, expectedMsgs)
		return nil
	}
	fetch := func(ctx context.Context) (kafka.Message, error) {
		return msg, nil
	}
	commit := func(ctx context.Context, msgs ...kafka.Message) error {
		cancel()
		return nil
	}
	r := &testFetchCommitter{
		fetch:  fetch,
		commit: commit,
	}
	c := &WaitConsumer{
		Producer: p,
	}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestWaitConsumerFetchContextDone(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	fetch := func(ctx context.Context) (kafka.Message, error) {
		cancel()
		return kafka.Message{}, ctx.Err()
	}
	r := &testFetchCommitter{
		fetch: fetch,
	}
	c := &WaitConsumer{}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestWaitConsumerWaitHeaderUndefined(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	msg := kafka.Message{
		Value: []byte("test"),
	}
	p := func(ctx context.Context, msgs ...kafka.Message) error {
		return nil
	}
	fetch := func(ctx context.Context) (kafka.Message, error) {
		return msg, nil
	}
	commit := func(ctx context.Context, msgs ...kafka.Message) error {
		cancel()
		return nil
	}
	r := &testFetchCommitter{
		fetch:  fetch,
		commit: commit,
	}
	c := &WaitConsumer{
		Producer: p,
	}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestWaitConsumerWaitHeaderInvalid(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	msg := kafka.Message{
		Value: []byte("test"),
		Headers: []kafka.Header{
			{
				Key:   WaitUntilHeader,
				Value: []byte("invalid"),
			},
		},
	}
	p := func(ctx context.Context, msgs ...kafka.Message) error {
		return nil
	}
	fetch := func(ctx context.Context) (kafka.Message, error) {
		return msg, nil
	}
	commit := func(ctx context.Context, msgs ...kafka.Message) error {
		cancel()
		return nil
	}
	r := &testFetchCommitter{
		fetch:  fetch,
		commit: commit,
	}
	c := &WaitConsumer{
		Producer: p,
	}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestWaitConsumerWaitHeaderPast(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	msg := kafka.Message{
		Value: []byte("test"),
		Headers: []kafka.Header{
			{
				Key:   WaitUntilHeader,
				Value: []byte(timeutils.Now().Add(-1 * time.Minute).Format(waitUntilHeaderLayout)),
			},
		},
	}
	p := func(ctx context.Context, msgs ...kafka.Message) error {
		return nil
	}
	fetch := func(ctx context.Context) (kafka.Message, error) {
		return msg, nil
	}
	commit := func(ctx context.Context, msgs ...kafka.Message) error {
		cancel()
		return nil
	}
	r := &testFetchCommitter{
		fetch:  fetch,
		commit: commit,
	}
	c := &WaitConsumer{
		Producer: p,
	}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestWaitConsumerWaitContextDone(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Microsecond)
	defer cancel()
	msg := kafka.Message{
		Value: []byte("test"),
		Headers: []kafka.Header{
			{
				Key:   WaitUntilHeader,
				Value: []byte(timeutils.Now().Add(1 * time.Minute).Format(waitUntilHeaderLayout)),
			},
		},
	}
	fetch := func(ctx context.Context) (kafka.Message, error) {
		return msg, nil
	}
	r := &testFetchCommitter{
		fetch: fetch,
	}
	c := &WaitConsumer{}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestWaitConsumerErrorFetch(t *testing.T) {
	ctx := context.Background()
	fetch := func(ctx context.Context) (kafka.Message, error) {
		return kafka.Message{}, errors.New("error")
	}
	r := &testFetchCommitter{
		fetch: fetch,
	}
	c := &WaitConsumer{}
	err := c.Consume(ctx, r)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestWaitConsumerErrorProducer(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	msg := kafka.Message{}
	p := func(ctx context.Context, msgs ...kafka.Message) error {
		cancel()
		return errors.New("error")
	}
	fetch := func(ctx context.Context) (kafka.Message, error) {
		return msg, nil
	}
	r := &testFetchCommitter{
		fetch: fetch,
	}
	c := &WaitConsumer{
		Producer: p,
	}
	err := c.Consume(ctx, r)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestWaitConsumerErrorCommit(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	msg := kafka.Message{}
	p := func(ctx context.Context, msgs ...kafka.Message) error {
		return nil
	}
	fetch := func(ctx context.Context) (kafka.Message, error) {
		return msg, nil
	}
	commit := func(ctx context.Context, msgs ...kafka.Message) error {
		cancel()
		return errors.New("error")
	}
	r := &testFetchCommitter{
		fetch:  fetch,
		commit: commit,
	}
	c := &WaitConsumer{
		Producer: p,
	}
	err := c.Consume(ctx, r)
	if err == nil {
		t.Fatal("no error")
	}
}
