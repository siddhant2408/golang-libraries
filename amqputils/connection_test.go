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

func TestConnectionManager(t *testing.T) {
	ctx := context.Background()
	com := amqptest.NewConnectionManager(t, testVhost)
	chn, err := com.Channel(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = chn.Close()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = com.Close()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestConnectionManagerErrorLock(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	com := &amqputils.ConnectionManager{}
	defer com.Close() //nolint:errcheck
	_, err := com.Channel(ctx)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestConnectionManagerErrorDial(t *testing.T) {
	ctx := context.Background()
	com := &amqputils.ConnectionManager{
		Dial: func(ctx context.Context) (*amqp.Connection, error) {
			return nil, errors.New("error")
		},
	}
	defer com.Close() //nolint:errcheck
	_, err := com.Channel(ctx)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestNewConnectionManagerURLs(t *testing.T) {
	ctx := context.Background()
	amqptest.CheckAvailable(t)
	u := amqptest.GetURLVHost(t, testVhost)
	com := amqputils.NewConnectionManagerURLs([]string{u})
	chn, err := com.Channel(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = chn.Close()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = com.Close()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
