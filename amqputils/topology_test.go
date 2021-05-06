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

const testVhost = "test"

func TestInitTopology(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	cg := amqputils.NewChannelGetterConnection(conn)
	err := amqputils.InitTopology(ctx, cg, amqputils.Topology{
		Exchanges: []amqputils.ExchangeConfig{
			{
				Name:       "X_1",
				Type:       amqp.ExchangeFanout,
				AutoDelete: true,
			},
			{
				Name:       "X_2",
				Type:       amqp.ExchangeFanout,
				AutoDelete: true,
				Bindings: []amqputils.ExchangeBinding{
					{
						Source: "X_1",
					},
				},
			},
		},
		Queues: []amqputils.QueueConfig{
			{
				Name:       "Q_1",
				AutoDelete: true,
				Bindings: []amqputils.QueueBinding{
					{
						Exchange: "X_2",
					},
				},
			},
		},
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestInitTopologyErrorGetChannel(t *testing.T) {
	cg := func(context.Context) (*amqp.Channel, error) {
		return nil, errors.New("error")
	}
	err := amqputils.InitTopology(context.Background(), cg, amqputils.Topology{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestInitTopologyErrorExchangeDeclare(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	cg := amqputils.NewChannelGetterConnection(conn)
	err := amqputils.InitTopology(ctx, cg, amqputils.Topology{
		Exchanges: []amqputils.ExchangeConfig{
			{
				Type: "invalid",
			},
		},
	})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestInitTopologyErrorExchangeBind(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	cg := amqputils.NewChannelGetterConnection(conn)
	err := amqputils.InitTopology(ctx, cg, amqputils.Topology{
		Exchanges: []amqputils.ExchangeConfig{
			{
				Name:       "X_3",
				Type:       amqp.ExchangeFanout,
				AutoDelete: true,
				Bindings: []amqputils.ExchangeBinding{
					{
						Source: "invalid",
					},
				},
			},
		},
	})
	if err == nil {
		t.Fatal("no error")
	}
	chn, err := cg(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	defer chn.Close() //nolint:errcheck
	err = chn.ExchangeDelete("X_3", false, false)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestInitTopologyErrorQueueDeclare(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	cg := amqputils.NewChannelGetterConnection(conn)
	err := amqputils.InitTopology(ctx, cg, amqputils.Topology{
		Queues: []amqputils.QueueConfig{
			{
				Name:       "amq.invalid",
				AutoDelete: true,
			},
		},
	})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestInitTopologyErrorQueueBind(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	cg := amqputils.NewChannelGetterConnection(conn)
	err := amqputils.InitTopology(ctx, cg, amqputils.Topology{
		Queues: []amqputils.QueueConfig{
			{
				Name:       "Q_2",
				AutoDelete: true,
				Bindings: []amqputils.QueueBinding{
					{
						Exchange: "invalid",
					},
				},
			},
		},
	})
	if err == nil {
		t.Fatal("no error")
	}
	chn, err := cg(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	defer chn.Close() //nolint:errcheck
	_, err = chn.QueueDelete("Q_2", false, false, false)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
