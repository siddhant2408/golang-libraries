package kafkautils

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/ctxutils"
	"github.com/siddhant2408/golang-libraries/errors"
)

// Accumulator accumulate messages.
type Accumulator struct {
	Size    int
	Timeout time.Duration
}

// Accumulate accumulate messages.
//
// When the maximum number of messages is reached, or the timeout is reached, or the context is canceled, it returns the accumulated messages.
// If the fetcher returns an error, this error is returned, with the accumulated messages.
func (a *Accumulator) Accumulate(ctx context.Context, f Fetcher) ([]kafka.Message, error) {
	var msgs []kafka.Message
	if a.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, a.Timeout)
		defer cancel()
	}
	for {
		msg, err := f(ctx)
		if err != nil {
			if ctxutils.IsDone(ctx) {
				return msgs, nil
			}
			return nil, errors.Wrap(err, "fetch")
		}
		msgs = append(msgs, msg)
		if a.Size > 0 && len(msgs) >= a.Size {
			return msgs, nil
		}
	}
}
