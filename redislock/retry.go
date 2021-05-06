package redislock

import (
	"time"
)

// RetryStrategy represents a retry strategy.
// The boolean indicates if a retry must be attempted.
// The duration indicates how long to wait before the next retry.
type RetryStrategy interface {
	Retry() (time.Duration, bool)
}

// NoRetry does not retry.
func NoRetry() RetryStrategy {
	return &noRetry{}
}

type noRetry struct{}

func (r *noRetry) Retry() (time.Duration, bool) {
	return 0, false
}

// LimitRetry retries up to a limit.
func LimitRetry(r RetryStrategy, limit int) RetryStrategy {
	return &limitRetry{
		r:     r,
		limit: limit,
	}
}

type limitRetry struct {
	r     RetryStrategy
	limit int
}

func (r *limitRetry) Retry() (time.Duration, bool) {
	if r.limit <= 0 {
		return 0, false
	}
	r.limit--
	return r.r.Retry()
}

// DelayRetry retries after a delay.
func DelayRetry(d time.Duration) RetryStrategy {
	return &delayRetry{
		d: d,
	}
}

type delayRetry struct {
	d time.Duration
}

func (r *delayRetry) Retry() (time.Duration, bool) {
	return r.d, true
}

// ExponentialRetry retries after an exponential delay.
// The delay is multiplied by two, up to the maximum value.
func ExponentialRetry(start, max time.Duration) RetryStrategy {
	return &exponentialRetry{
		d:   start,
		max: max,
	}
}

type exponentialRetry struct {
	d   time.Duration
	max time.Duration
}

func (r *exponentialRetry) Retry() (time.Duration, bool) {
	if r.d > r.max {
		return r.max, true
	}
	d := r.d
	r.d *= 2
	return d, true
}
