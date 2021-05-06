package kafkautils

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestAccumulator(t *testing.T) {
	ctx := context.Background()
	f := func(context.Context) (kafka.Message, error) {
		return kafka.Message{}, nil
	}
	a := Accumulator{
		Size: 10,
	}
	msgs, err := a.Accumulate(ctx, f)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if len(msgs) != a.Size {
		t.Fatalf("unexpected message count: got %d, want %d", len(msgs), a.Size)
	}
}

func TestAccumulatorTimeout(t *testing.T) {
	ctx := context.Background()
	f := func(ctx context.Context) (kafka.Message, error) {
		<-ctx.Done()
		return kafka.Message{}, ctx.Err()
	}
	a := Accumulator{
		Timeout: 1 * time.Nanosecond,
	}
	_, err := a.Accumulate(ctx, f)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestAccumulatorContextCancel(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	f := func(ctx context.Context) (kafka.Message, error) {
		cancel()
		<-ctx.Done()
		return kafka.Message{}, ctx.Err()
	}
	a := Accumulator{}
	_, err := a.Accumulate(ctx, f)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestAccumulatorErrorFetcher(t *testing.T) {
	ctx := context.Background()
	f := func(context.Context) (kafka.Message, error) {
		return kafka.Message{}, errors.New("error")
	}
	a := Accumulator{}
	_, err := a.Accumulate(ctx, f)
	if err == nil {
		t.Fatal("no error")
	}
}
