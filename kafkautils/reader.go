package kafkautils

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/ctxutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/goroutine"
)

// ConsumeReader consumes a Reader.
func ConsumeReader(ctx context.Context, cfg kafka.ReaderConfig, c ReaderConsumer, errFunc func(context.Context, error)) {
	for !ctxutils.IsDone(ctx) {
		func() {
			r := kafka.NewReader(cfg)
			defer func() {
				err := r.Close()
				if err != nil {
					err = errors.Wrap(err, "close Kafka Reader")
					errFunc(ctx, err)
				}
			}()
			err := c(ctx, r)
			if err != nil {
				err = errors.Wrap(err, "Kafka Reader")
				errFunc(ctx, err)
			}
		}()
	}
}

// ConsumeReaders calls ConsumeReader concurrently.
func ConsumeReaders(ctx context.Context, cfg kafka.ReaderConfig, count int, c ReaderConsumer, errFunc func(context.Context, error)) {
	goroutine.RunN(count, func() {
		ConsumeReader(ctx, cfg, c, errFunc)
	})
}

// FetchCommitter fetches and commits messages.
// It is implemented by Reader.
type FetchCommitter interface {
	FetchMessage(context.Context) (kafka.Message, error)
	CommitMessages(context.Context, ...kafka.Message) error
}

// Fetcher fetches messages.
type Fetcher func(context.Context) (kafka.Message, error)

// ReaderConsumer is a consumer for a Reader.
type ReaderConsumer func(context.Context, FetchCommitter) error
