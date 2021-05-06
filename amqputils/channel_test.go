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

func TestChannelPool(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	cp := &amqputils.ChannelPool{
		Channel: amqputils.NewChannelGetterConnection(conn),
	}
	for i := 0; i < 5; i++ {
		chn, err := cp.Get(ctx)
		if err != nil {
			testutils.FatalErr(t, err)
		}
		cp.Put(chn)
	}
	err := cp.Close()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestChannelPoolErrorLock(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	cp := &amqputils.ChannelPool{}
	_, err := cp.Get(ctx)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestChannelPoolErrorOpen(t *testing.T) {
	ctx := context.Background()
	cp := &amqputils.ChannelPool{
		Channel: func(ctx context.Context) (*amqp.Channel, error) {
			return nil, errors.New("error")
		},
	}
	_, err := cp.Get(ctx)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestChannelPoolRun(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	cp := &amqputils.ChannelPool{
		Channel: amqputils.NewChannelGetterConnection(conn),
	}
	err := cp.Run(ctx, func(ctx context.Context, ch *amqp.Channel) error {
		return nil
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = cp.Close()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestChannelPoolRunErrorGetChannel(t *testing.T) {
	ctx := context.Background()
	cp := &amqputils.ChannelPool{
		Channel: func(ctx context.Context) (*amqp.Channel, error) {
			return nil, errors.New("error")
		},
	}
	err := cp.Run(ctx, func(ctx context.Context, ch *amqp.Channel) error {
		t.Fatal("should not be called")
		return nil
	})
	if err == nil {
		t.Fatal("no error")
	}
	err = cp.Close()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestChannelPoolRunErrorRecycle(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	cp := &amqputils.ChannelPool{
		Channel: amqputils.NewChannelGetterConnection(conn),
	}
	err := cp.Run(ctx, func(ctx context.Context, ch *amqp.Channel) error {
		return errors.New("error")
	})
	if err == nil {
		t.Fatal("no error")
	}
	err = cp.Close()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestChannelPoolRunErrorDiscard(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	cp := &amqputils.ChannelPool{
		Channel: amqputils.NewChannelGetterConnection(conn),
	}
	err := cp.Run(ctx, func(ctx context.Context, ch *amqp.Channel) error {
		return amqp.ErrClosed
	})
	if err == nil {
		t.Fatal("no error")
	}
	err = cp.Close()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestChannelPoolErrorClose(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	cp := &amqputils.ChannelPool{
		Channel: amqputils.NewChannelGetterConnection(conn),
	}
	chn, err := cp.Get(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	cp.Put(chn)
	_ = conn.Close()
	err = cp.Close()
	if err == nil {
		t.Fatal("no error")
	}
}
