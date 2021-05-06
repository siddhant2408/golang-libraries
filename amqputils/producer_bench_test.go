package amqputils_test

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/amqptest"
	"github.com/siddhant2408/golang-libraries/amqputils"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/streadway/amqp"
)

func BenchmarkSimpleProducerNoConfirm(b *testing.B) {
	conn := amqptest.NewConnection(b, testVhost)
	p := &amqputils.SimpleProducer{
		Channel: amqputils.NewChannelGetterConnection(conn),
		Confirm: false,
	}
	m := amqp.Publishing{
		Body: []byte("test"),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := p.Produce(context.Background(), "", "_test", false, false, m)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
	err := p.Close()
	if err != nil {
		testutils.FatalErr(b, err)
	}
}

func BenchmarkSimpleProducerConfirm(b *testing.B) {
	conn := amqptest.NewConnection(b, testVhost)
	p := &amqputils.SimpleProducer{
		Channel: amqputils.NewChannelGetterConnection(conn),
		Confirm: true,
	}
	m := amqp.Publishing{
		Body: []byte("test"),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := p.Produce(context.Background(), "", "_test", false, false, m)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
	err := p.Close()
	if err != nil {
		testutils.FatalErr(b, err)
	}
}
