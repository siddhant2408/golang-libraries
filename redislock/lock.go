package redislock

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/siddhant2408/golang-libraries/ctxutils"
	"github.com/siddhant2408/golang-libraries/errors"
)

// Locker manages the locks.
type Locker struct {
	Client interface {
		SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
		Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
		EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd
		ScriptExists(ctx context.Context, scripts ...string) *redis.BoolSliceCmd
		ScriptLoad(ctx context.Context, script string) *redis.StringCmd
	}
}

// Lock attempts to acquire a lock.
// It may return ErrNotObtained if the lock could not be obtained.
func (lr *Locker) Lock(ctx context.Context, opts ...Option) (_ *Lock, err error) {
	span, spanFinish := startTraceSpan(&ctx, "lock", &err)
	defer spanFinish()
	o, err := getOptions(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "get options")
	}
	setTraceSpanTag(span, "key", o.key)
	setTraceSpanTag(span, "value", o.value)
	setTraceSpanTag(span, "ttl", o.ttl.String())
	l, err := lr.lock(ctx, o)
	if err != nil {
		err = wrapErrorValue(err, "key", o.key)
		err = wrapErrorValue(err, "value", o.value)
		err = wrapErrorValue(err, "ttl", o.ttl)
		return nil, err
	}
	return l, nil
}

func (lr *Locker) lock(ctx context.Context, o *options) (*Lock, error) {
	for attempt := 0; ; attempt++ {
		l, err := lr.lockAttempt(ctx, o, attempt)
		if err != nil {
			err = errors.Wrap(err, "attempt")
			err = wrapErrorValue(err, "attempt", attempt)
			return nil, err
		}
		if l != nil {
			return l, nil
		}
	}
}

func (lr *Locker) lockAttempt(ctx context.Context, o *options, attempt int) (_ *Lock, err error) {
	span, spanFinish := startTraceSpan(&ctx, "lock.attempt", &err)
	defer spanFinish()
	setTraceSpanTag(span, "attempt", attempt)
	ok, err := lr.lockRedis(ctx, o.key, o.value, o.ttl)
	if err != nil {
		return nil, errors.Wrap(err, "Redis")
	}
	if ok {
		return &Lock{
			Key:   o.key,
			Value: o.value,
		}, nil
	}
	err = lr.lockRetry(ctx, o.retry)
	if err != nil {
		return nil, errors.Wrap(err, "retry")
	}
	return nil, nil
}

func (lr *Locker) lockRedis(ctx context.Context, key string, value string, ttl time.Duration) (_ bool, err error) {
	_, spanFinish := startTraceSpan(&ctx, "lock.redis", &err)
	defer spanFinish()
	ok, err := lr.Client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return false, errors.Wrap(err, "SETNX")
	}
	return ok, nil
}

func (lr *Locker) lockRetry(ctx context.Context, r RetryStrategy) (err error) {
	span, spanFinish := startTraceSpan(&ctx, "lock.retry", &err)
	defer spanFinish()
	if ctxutils.IsDone(ctx) {
		return errors.Wrap(ctx.Err(), "")
	}
	d, ok := r.Retry()
	if !ok {
		return errors.Wrap(ErrNotObtained, "")
	}
	setTraceSpanTag(span, "retry.wait", d.String())
	if d <= 0 {
		return nil
	}
	tm := time.NewTimer(d)
	defer tm.Stop()
	select {
	case <-tm.C:
		return nil
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "")
	}
}

var scriptUnlock = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)

// Unlock releases the lock.
// It may return ErrNotHeld if the lock is not held.
func (lr *Locker) Unlock(ctx context.Context, l *Lock) (err error) {
	span, spanFinish := startTraceSpan(&ctx, "unlock", &err)
	defer spanFinish()
	setTraceSpanTag(span, "key", l.Key)
	setTraceSpanTag(span, "value", l.Value)
	err = lr.unlock(ctx, l)
	if err != nil {
		err = wrapErrorValue(err, "key", l.Key)
		err = wrapErrorValue(err, "value", l.Value)
		return err
	}
	return nil
}

func (lr *Locker) unlock(ctx context.Context, l *Lock) error {
	res, err := scriptUnlock.Run(ctx, lr.Client, []string{l.Key}, l.Value).Int64()
	if err != nil {
		return errors.Wrap(err, "script")
	}
	if res != 1 {
		return errors.Wrap(ErrNotHeld, "")
	}
	return nil
}

var scriptRefresh = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pexpire", KEYS[1], ARGV[2]) else return 0 end`)

// Refresh refreshes the TTL of a lock.
// It may return ErrNotHeld if the lock is not held.
func (lr *Locker) Refresh(ctx context.Context, l *Lock, ttl time.Duration) (err error) {
	span, spanFinish := startTraceSpan(&ctx, "refresh", &err)
	defer spanFinish()
	setTraceSpanTag(span, "key", l.Key)
	setTraceSpanTag(span, "value", l.Value)
	setTraceSpanTag(span, "ttl", ttl.String())
	err = lr.refresh(ctx, l, ttl)
	if err != nil {
		err = wrapErrorValue(err, "key", l.Key)
		err = wrapErrorValue(err, "value", l.Value)
		err = wrapErrorValue(err, "ttl", ttl)
		return err
	}
	return nil
}

func (lr *Locker) refresh(ctx context.Context, l *Lock, ttl time.Duration) error {
	ttlStr := strconv.FormatInt(ttl.Milliseconds(), 10)
	res, err := scriptRefresh.Run(ctx, lr.Client, []string{l.Key}, l.Value, ttlStr).Int64()
	if err != nil {
		return errors.Wrap(err, "script")
	}
	if res != 1 {
		return errors.Wrap(ErrNotHeld, "")
	}
	return nil
}

var scriptTTL = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pttl", KEYS[1]) else return -3 end`)

// TTL returns the TTL for a lock.
// If the key has no TTL, it returns a negative value.
// It may return ErrNotHeld if the lock is not held.
func (lr *Locker) TTL(ctx context.Context, l *Lock) (ttl time.Duration, err error) {
	span, spanFinish := startTraceSpan(&ctx, "ttl", &err)
	defer spanFinish()
	setTraceSpanTag(span, "key", l.Key)
	setTraceSpanTag(span, "value", l.Value)
	ttl, err = lr.ttl(ctx, l)
	if err != nil {
		err = wrapErrorValue(err, "key", l.Key)
		err = wrapErrorValue(err, "value", l.Value)
		return 0, err
	}
	return ttl, nil
}

func (lr *Locker) ttl(ctx context.Context, l *Lock) (time.Duration, error) {
	ttlMS, err := scriptTTL.Run(ctx, lr.Client, []string{l.Key}, l.Value).Int64()
	if err != nil {
		return 0, errors.Wrap(err, "script")
	}
	if ttlMS == -3 {
		return 0, errors.Wrap(ErrNotHeld, "")
	}
	ttl := time.Duration(ttlMS) * time.Millisecond
	return ttl, nil
}

// Lock represent a lock value.
type Lock struct {
	Key   string
	Value string
}
