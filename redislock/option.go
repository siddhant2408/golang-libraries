package redislock

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/siddhant2408/golang-libraries/errors"
)

type options struct {
	key   string
	value string
	ttl   time.Duration
	retry RetryStrategy
}

func getOptions(opts ...Option) (*options, error) {
	o := newOptions()
	for _, opt := range opts {
		opt(o)
	}
	if o.value == "" {
		v, err := generateRandomValue()
		if err != nil {
			return nil, errors.Wrap(err, "generate random value")
		}
		o.value = v
	}
	return o, nil
}

func newOptions() *options {
	return &options{
		key:   "lock",
		ttl:   5 * time.Minute,
		retry: NoRetry(),
	}
}

// Option represents an option.
type Option func(*options)

// Key defines the key.
// The default value is "lock".
func Key(key string) Option {
	return func(o *options) {
		o.key = key
	}
}

// Value defines the value.
// The default value is 16 random bytes hex encoded.
func Value(value string) Option {
	return func(o *options) {
		o.value = value
	}
}

// TTL defines the TTL of the key.
// The default value is 5 minutes.
func TTL(ttl time.Duration) Option {
	return func(o *options) {
		o.ttl = ttl
	}
}

// Retry defines the retry strategy.
// The default value is no retry.
func Retry(r RetryStrategy) Option {
	return func(o *options) {
		o.retry = r
	}
}

const randomValueByteCount = 16

func generateRandomValue() (string, error) {
	var buf [randomValueByteCount]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		return "", errors.Wrap(err, "rand read")
	}
	s := hex.EncodeToString(buf[:])
	return s, nil
}
