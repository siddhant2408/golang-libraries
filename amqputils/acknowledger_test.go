package amqputils_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/siddhant2408/golang-libraries/amqputils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestAck(t *testing.T) {
	msg := &testDelivery{
		t:           t,
		expectedAck: true,
	}
	err := amqputils.Ack.Acknowledge(msg)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestNackRequeue(t *testing.T) {
	msg := &testDelivery{
		t:               t,
		expectedNack:    true,
		expectedRequeue: true,
	}
	err := amqputils.NackRequeue.Acknowledge(msg)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestNackDiscard(t *testing.T) {
	msg := &testDelivery{
		t:               t,
		expectedNack:    true,
		expectedRequeue: false,
	}
	err := amqputils.NackDiscard.Acknowledge(msg)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestAcknowledgerError(t *testing.T) {
	err := errors.New("test")
	err = amqputils.ErrorWithAcknowledger(err, amqputils.Ack)
	ack := amqputils.GetErrorAcknowledger(err)
	if ack != amqputils.Ack {
		t.Fatalf("unexpected Acknowledger: got %v, want %v", ack, amqputils.Ack)
	}
	_ = fmt.Sprintf("%s %+v", err, err)
}

func TestErrorWithAcknowledgerNil(t *testing.T) {
	err := amqputils.ErrorWithAcknowledger(nil, amqputils.Ack)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestGetErrorAcknowledgerNoMatch(t *testing.T) {
	err := errors.New("test")
	ack := amqputils.GetErrorAcknowledger(err)
	if ack != nil {
		t.Fatalf("unexpected Acknowledger: got %v, want %v", ack, nil)
	}
}

func BenchmarkAcknowledgerErrorFormat(b *testing.B) {
	err := errors.New("error")
	err = amqputils.ErrorWithAcknowledger(err, amqputils.NackDiscard)
	for i := 0; i < b.N; i++ {
		fmt.Fprint(io.Discard, err)
	}
}

type testDelivery struct {
	t               *testing.T
	expectedAck     bool
	expectedNack    bool
	expectedRequeue bool
}

func (d *testDelivery) Ack(multiple bool) error {
	if !d.expectedAck {
		d.t.Fatal("unexpected call to Ack")
	}
	return nil
}

func (d *testDelivery) Nack(multiple bool, requeue bool) error {
	if !d.expectedNack {
		d.t.Fatal("unexpected call to Nack")
	}
	if requeue != d.expectedRequeue {
		d.t.Fatalf("unexpected requeue: got %t, want %t", requeue, d.expectedRequeue)
	}
	return nil
}

type testAMQPAcknowledger struct {
	testAMQPAcknowledgerAck
	testAMQPAcknowledgerNack
	testAMQPAcknowledgerReject
}

type testAMQPAcknowledgerAck func(tag uint64, multiple bool) error

func (f testAMQPAcknowledgerAck) Ack(tag uint64, multiple bool) error {
	return f(tag, multiple)
}

type testAMQPAcknowledgerNack func(tag uint64, multiple bool, requeue bool) error

func (f testAMQPAcknowledgerNack) Nack(tag uint64, multiple bool, requeue bool) error {
	return f(tag, multiple, requeue)
}

type testAMQPAcknowledgerReject func(tag uint64, multiple bool) error

func (f testAMQPAcknowledgerReject) Reject(tag uint64, multiple bool) error {
	return f(tag, multiple)
}
