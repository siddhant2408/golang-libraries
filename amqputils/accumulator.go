package amqputils

import (
	"context"
	"time"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/streadway/amqp"
)

// Accumulator accumulates messages from a channel.
type Accumulator struct {
	Size  int
	Delay time.Duration
}

// Accumulate accumulates messages.
//
// When the maximum number of messages is reached, or the delay is reached, or the context is canceled, it returns the accumulated messages.
// If the channel is closed, an error is returned, with the accumulated messages.
func (a *Accumulator) Accumulate(ctx context.Context, ch <-chan amqp.Delivery) ([]amqp.Delivery, error) {
	var dlvs []amqp.Delivery
	tm := time.NewTimer(a.Delay)
	defer tm.Stop()
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return dlvs, errors.New("channel closed")
			}
			dlvs = append(dlvs, msg)
			if len(dlvs) >= a.Size {
				return dlvs, nil
			}
		case <-tm.C:
			return dlvs, nil
		case <-ctx.Done():
			return dlvs, nil
		}
	}
}
