package amqputils

import (
	"fmt"

	"github.com/siddhant2408/golang-libraries/errors"
)

// Acknowledger acknowledges messages.
type Acknowledger interface {
	Acknowledge(Delivery) error
	String() string
}

// Ack is an Acknowledger that acknowledges the message.
const Ack = ack("ack")

type ack string

func (a ack) Acknowledge(msg Delivery) error {
	return errors.Wrap(msg.Ack(false), string(a))
}

func (a ack) String() string {
	return string(a)
}

// NackRequeue is an Acknowledger that negatively acknowledges the message with requeue=true.
const NackRequeue = nackRequeue("nack requeue")

type nackRequeue string

func (a nackRequeue) Acknowledge(msg Delivery) error {
	return errors.Wrap(msg.Nack(false, true), string(a))
}

func (a nackRequeue) String() string {
	return string(a)
}

// NackDiscard is an Acknowledger that negatively acknowledges the message with requeue=false.
const NackDiscard = nackDiscard("nack discard")

type nackDiscard string

func (a nackDiscard) Acknowledge(msg Delivery) error {
	return errors.Wrap(msg.Nack(false, false), string(a))
}

func (a nackDiscard) String() string {
	return string(a)
}

// Delivery represents an AMQP delivery.
type Delivery interface {
	Ack(multiple bool) error
	Nack(multiple bool, requeue bool) error
}

// ErrorWithAcknowledger adds an Acknowledger to the error.
// It is used by Consumer.
func ErrorWithAcknowledger(err error, ack Acknowledger) error {
	if err == nil {
		return nil
	}
	return &acknowledgerError{
		err: err,
		ack: ack,
	}
}

type acknowledgerError struct {
	err error
	ack Acknowledger
}

func (err *acknowledgerError) AMQPAcknowledger() Acknowledger {
	return err.ack
}

func (err *acknowledgerError) WriteErrorMessage(w errors.Writer, verbose bool) bool {
	_, _ = w.WriteString("AMQP ")
	_, _ = w.WriteString(err.ack.String())
	return true
}

func (err *acknowledgerError) Error() string                 { return errors.Error(err) }
func (err *acknowledgerError) Format(s fmt.State, verb rune) { errors.Format(err, s, verb) }
func (err *acknowledgerError) Unwrap() error                 { return err.err }

// GetErrorAcknowledger returns the Acknowledger associated to the error.
func GetErrorAcknowledger(err error) Acknowledger {
	var werr *acknowledgerError
	ok := errors.As(err, &werr)
	if ok {
		return werr.AMQPAcknowledger()
	}
	return nil
}
