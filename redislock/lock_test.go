package redislock

import (
	"context"
	"testing"
	"time"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/redistest"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestLocker(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l, err := lr.Lock(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = lr.Unlock(ctx, l)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	exists, err := client.Exists(ctx, l.Key).Result()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if exists != 0 {
		t.Fatalf("unexpected exists: got %d, want %d", exists, 0)
	}
}

func TestLockerLockKey(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l, err := lr.Lock(ctx, Key("test"))
	if err != nil {
		testutils.FatalErr(t, err)
	}
	exists, err := client.Exists(ctx, "test").Result()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if exists != 1 {
		t.Fatalf("unexpected exists: got %d, want %d", exists, 1)
	}
	err = lr.Unlock(ctx, l)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestLockerLockValue(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l, err := lr.Lock(ctx, Value("test"))
	if err != nil {
		testutils.FatalErr(t, err)
	}
	value, err := client.Get(ctx, l.Key).Result()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if value != "test" {
		t.Fatalf("unexpected value: got %s, want %s", value, "test")
	}
	err = lr.Unlock(ctx, l)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestLockerLockTTL(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l, err := lr.Lock(ctx, TTL(1*time.Minute))
	if err != nil {
		testutils.FatalErr(t, err)
	}
	ttl, err := client.TTL(ctx, l.Key).Result()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	d := 1*time.Minute - ttl
	if d > 5*time.Second {
		t.Fatalf("unexpected TTL: got %v, want %v", ttl, 1*time.Minute)
	}
	err = lr.Unlock(ctx, l)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestLockerErrorRetryNotObtained(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l, err := lr.Lock(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	_, err = lr.Lock(ctx, Retry(LimitRetry(DelayRetry(1*time.Nanosecond), 3)))
	if err == nil {
		t.Fatal("no error")
	}
	if !errors.Is(err, ErrNotObtained) {
		t.Fatalf("unexpected error: got %v, want %v", err, ErrNotObtained)
	}
	err = lr.Unlock(ctx, l)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestLockerErrorRetryNotObtainedZero(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l, err := lr.Lock(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	_, err = lr.Lock(ctx, Retry(LimitRetry(DelayRetry(0), 3)))
	if err == nil {
		t.Fatal("no error")
	}
	if !errors.Is(err, ErrNotObtained) {
		t.Fatalf("unexpected error: got %v, want %v", err, ErrNotObtained)
	}
	err = lr.Unlock(ctx, l)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestLockerErrorRetryContextCanceled(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l, err := lr.Lock(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	ctxCanceled, cancel := context.WithCancel(ctx)
	cancel()
	_, err = lr.Lock(ctxCanceled)
	if err == nil {
		t.Fatal("no error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("unexpected error: got %v, want %v", err, context.Canceled)
	}
	err = lr.Unlock(ctx, l)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestLockerUnlockErrorNotHeldNotLocked(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l := &Lock{
		Key:   "test",
		Value: "test",
	}
	err := lr.Unlock(ctx, l)
	if err == nil {
		t.Fatal("no error")
	}
	if !errors.Is(err, ErrNotHeld) {
		t.Fatalf("unexpected error: got %v, want %v", err, ErrNotHeld)
	}
}

func TestLockerUnlockErrorNotHeldAlreadyLocked(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l1, err := lr.Lock(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	l2 := &Lock{
		Key:   l1.Key,
		Value: "test",
	}
	err = lr.Unlock(ctx, l2)
	if err == nil {
		t.Fatal("no error")
	}
	if !errors.Is(err, ErrNotHeld) {
		t.Fatalf("unexpected error: got %v, want %v", err, ErrNotHeld)
	}
	err = lr.Unlock(ctx, l1)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestLockerRefresh(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l, err := lr.Lock(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = lr.Refresh(ctx, l, 10*time.Minute)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	ttl, err := client.TTL(ctx, l.Key).Result()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	d := 10*time.Minute - ttl
	if d > 5*time.Second {
		t.Fatalf("unexpected TTL: got %v, want %v", ttl, 10*time.Minute)
	}
	err = lr.Unlock(ctx, l)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestLockerRefreshErrorNotHeldNotLocked(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l := &Lock{
		Key:   "test",
		Value: "test",
	}
	err := lr.Refresh(ctx, l, 10*time.Minute)
	if err == nil {
		t.Fatal("no error")
	}
	if !errors.Is(err, ErrNotHeld) {
		t.Fatalf("unexpected error: got %v, want %v", err, ErrNotHeld)
	}
}

func TestLockerRefreshErrorNotHeldAlreadyLocked(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l1, err := lr.Lock(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	l2 := &Lock{
		Key:   l1.Key,
		Value: "test",
	}
	err = lr.Refresh(ctx, l2, 10*time.Minute)
	if err == nil {
		t.Fatal("no error")
	}
	if !errors.Is(err, ErrNotHeld) {
		t.Fatalf("unexpected error: got %v, want %v", err, ErrNotHeld)
	}
	err = lr.Unlock(ctx, l1)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestLockerTTL(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l, err := lr.Lock(ctx, TTL(1*time.Minute))
	if err != nil {
		testutils.FatalErr(t, err)
	}
	ttl, err := lr.TTL(ctx, l)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	d := 1*time.Minute - ttl
	if d > 5*time.Second {
		t.Fatalf("unexpected TTL: got %v, want %v", ttl, 1*time.Minute)
	}
	err = lr.Unlock(ctx, l)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestLockerTTLErrorNotHeldNotLocked(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l := &Lock{
		Key:   "test",
		Value: "test",
	}
	_, err := lr.TTL(ctx, l)
	if err == nil {
		t.Fatal("no error")
	}
	if !errors.Is(err, ErrNotHeld) {
		t.Fatalf("unexpected error: got %v, want %v", err, ErrNotHeld)
	}
}

func TestLockerTTLErrorNotHeldAlreadyLocked(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l1, err := lr.Lock(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	l2 := &Lock{
		Key:   l1.Key,
		Value: "test",
	}
	_, err = lr.TTL(ctx, l2)
	if err == nil {
		t.Fatal("no error")
	}
	if !errors.Is(err, ErrNotHeld) {
		t.Fatalf("unexpected error: got %v, want %v", err, ErrNotHeld)
	}
	err = lr.Unlock(ctx, l1)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestLockerTTLFoo(t *testing.T) {
	ctx := context.Background()
	client := redistest.NewClient(t)
	lr := &Locker{
		Client: client,
	}
	l, err := lr.Lock(ctx, TTL(0))
	if err != nil {
		testutils.FatalErr(t, err)
	}
	ttl, err := lr.TTL(ctx, l)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if ttl != -1*time.Millisecond {
		t.Fatalf("unexpected TTL: got %v, want %v", ttl, -1*time.Millisecond)
	}
	err = lr.Unlock(ctx, l)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
