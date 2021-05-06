package amqputils

import (
	"context"
	"strconv"
	"time"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/streadway/amqp"
)

const (
	retryHeaderAttempts = "retry-attempts"
)

// Retryer allows to retry messages.
type Retryer struct {
	Producer Producer
	Max      int64
	MaxAck   bool
	Delay    time.Duration
	Exchange string
	Key      string
}

// Retry retries a message.
//
// It stores the number of attempts in the "retry-attempts" header.
// If the number of attempts is greater than or equal to the maximum value, it returns an error that discards the message.
// The number of attempts is incremented by 1 for each attempt.
//
// If Max is less than or equal to 0, there is not maximum attempts.
//
// MaxAck controls the behavior of the error that is returned if the maximum attempts is reached.
// If false it discards the message.
// If true it acknowledges the message.
//
// It produces the message to the given exchange and key.
// The message expiration is set with the delay.
//
// In case of success, it always returns an error that is ignored and acknowledge the message.
func (r *Retryer) Retry(ctx context.Context, dlv amqp.Delivery) error {
	at := r.getAttempts(dlv)
	if r.Max > 0 && at >= r.Max {
		err := errors.Newf("max retry reached: %d", r.Max)
		a := r.getMaxAcknowledger()
		err = ErrorWithAcknowledger(err, a)
		return err
	}
	pbl := deliveryToPublishing(dlv)
	if r.Max > 0 {
		at++
		pbl.Headers[retryHeaderAttempts] = at
	} else {
		delete(pbl.Headers, retryHeaderAttempts)
	}
	pbl.Expiration = strconv.FormatInt(int64(r.Delay/time.Millisecond), 10)
	err := r.Producer(ctx, r.Exchange, r.Key, false, false, pbl)
	if err != nil {
		return errors.Wrap(err, "produce")
	}
	err = errors.New("retry")
	err = errors.Ignore(err)
	err = ErrorWithAcknowledger(err, Ack)
	return err
}

func (r *Retryer) getAttempts(dlv amqp.Delivery) int64 {
	h, ok := dlv.Headers[retryHeaderAttempts]
	if !ok {
		return 0
	}
	a, ok := h.(int64)
	if !ok {
		return 0
	}
	return a
}

func (r *Retryer) getMaxAcknowledger() Acknowledger {
	if r.MaxAck {
		return Ack
	}
	return NackDiscard
}
