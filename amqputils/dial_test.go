package amqputils_test

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/amqptest"
	"github.com/siddhant2408/golang-libraries/amqputils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/streadway/amqp"
)

func TestDial(t *testing.T) {
	ctx := context.Background()
	amqptest.CheckAvailable(t)
	u := amqptest.GetURLVHost(t, testVhost)
	conn, err := amqputils.Dial(ctx, u)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = conn.Close()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestDialError(t *testing.T) {
	ctx := context.Background()
	u := "amqp://invalid:5672"
	_, err := amqputils.Dial(ctx, u)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestDialConfig(t *testing.T) {
	ctx := context.Background()
	amqptest.CheckAvailable(t)
	u := amqptest.GetURLVHost(t, testVhost)
	cfg := amqp.Config{}
	conn, err := amqputils.DialConfig(ctx, u, cfg)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = conn.Close()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestDialConfigError(t *testing.T) {
	ctx := context.Background()
	u := "amqp://invalid:5672"
	cfg := amqp.Config{}
	_, err := amqputils.DialConfig(ctx, u, cfg)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestURLsDialer(t *testing.T) {
	amqptest.CheckAvailable(t)
	ctx := context.Background()
	u := amqptest.GetURLVHost(t, testVhost)
	d := &amqputils.URLsDialer{
		URLs:    []string{u},
		DialURL: amqputils.Dial,
	}
	conn, err := d.Dial(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = conn.Close()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestURLsDialerErrorContextDone(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	u := amqptest.GetURLVHost(t, testVhost)
	d := &amqputils.URLsDialer{
		URLs: []string{u},
	}
	_, err := d.Dial(ctx)
	if err == nil {
		t.Fatal("no error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("unexpected error: got %v, want %v", err, context.Canceled)
	}
}

func TestURLsDialerErrorNoURL(t *testing.T) {
	ctx := context.Background()
	d := &amqputils.URLsDialer{}
	_, err := d.Dial(ctx)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestURLsDialerErrorAllURL(t *testing.T) {
	ctx := context.Background()
	u := "amqp://invalid:5672"
	d := &amqputils.URLsDialer{
		URLs: []string{u},
		DialURL: func(context.Context, string) (*amqp.Connection, error) {
			return nil, errors.New("error")
		},
	}
	_, err := d.Dial(ctx)
	if err == nil {
		t.Fatal("no error")
	}
}
