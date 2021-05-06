package amqputils_test

import (
	"context"
	"testing"
	"time"

	"github.com/siddhant2408/golang-libraries/amqputils"
	"github.com/siddhant2408/golang-libraries/goroutine"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/streadway/amqp"
)

func TestAccumulatorSize(t *testing.T) {
	size := 10
	ch := make(chan amqp.Delivery, size)
	for i := 0; i < size; i++ {
		ch <- amqp.Delivery{}
	}
	a := &amqputils.Accumulator{
		Size:  size,
		Delay: 10 * time.Second,
	}
	dlvs, err := a.Accumulate(context.Background(), ch)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if len(dlvs) != size {
		t.Fatalf("unexpected length: got %d, want %d", len(dlvs), size)
	}
}

func TestAccumulatorDelay(t *testing.T) {
	if testing.Short() {
		t.Skip("this test is slow")
	}
	ch := make(chan amqp.Delivery, 1)
	ch <- amqp.Delivery{}
	a := &amqputils.Accumulator{
		Size:  10,
		Delay: 1 * time.Second, // We need a long delay, otherwise the test can be flakky.
	}
	dlvs, err := a.Accumulate(context.Background(), ch)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if len(dlvs) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(dlvs), 1)
	}
}

func TestAccumulatorCancel(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ch := make(chan amqp.Delivery)
	waitSendDelivery := goroutine.Go(func() {
		ch <- amqp.Delivery{}
		cancel()
	})
	a := &amqputils.Accumulator{
		Size:  10,
		Delay: 10 * time.Second,
	}
	dlvs, err := a.Accumulate(ctx, ch)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	waitSendDelivery()
	if len(dlvs) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(dlvs), 1)
	}
}

func TestAccumulatorErrorChannelClosed(t *testing.T) {
	ch := make(chan amqp.Delivery)
	close(ch)
	a := &amqputils.Accumulator{
		Size:  10,
		Delay: 10 * time.Second,
	}
	_, err := a.Accumulate(context.Background(), ch)
	if err == nil {
		t.Fatal("no error")
	}
}
