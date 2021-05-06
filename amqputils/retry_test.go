package amqputils

import (
	"context"
	"testing"
	"time"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/streadway/amqp"
)

func TestRetryer(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		name        string
		dlv         amqp.Delivery
		expectedPbl amqp.Publishing
	}{
		{
			name: "NoAttempt",
			dlv: amqp.Delivery{
				Body: []byte("test"),
			},
			expectedPbl: amqp.Publishing{
				Headers: amqp.Table{
					retryHeaderAttempts: int64(1),
				},
				Expiration: "60000",
				Body:       []byte("test"),
			},
		},
		{
			name: "Attempt",
			dlv: amqp.Delivery{
				Headers: amqp.Table{
					retryHeaderAttempts: int64(5),
				},
				Body: []byte("test"),
			},
			expectedPbl: amqp.Publishing{
				Headers: amqp.Table{
					retryHeaderAttempts: int64(6),
				},
				Expiration: "60000",
				Body:       []byte("test"),
			},
		},
		{
			name: "HeaderWrongType",
			dlv: amqp.Delivery{
				Headers: amqp.Table{
					retryHeaderAttempts: "invalid",
				},
				Body: []byte("test"),
			},
			expectedPbl: amqp.Publishing{
				Headers: amqp.Table{
					retryHeaderAttempts: int64(1),
				},
				Expiration: "60000",
				Body:       []byte("test"),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := &Retryer{
				Max:      10,
				Delay:    1 * time.Minute,
				Exchange: "exchange",
				Key:      "key",
			}
			var pCalled testutils.CallCounter
			r.Producer = func(ctx context.Context, exchange, key string, mandatory, immediate bool, pbl amqp.Publishing) error {
				pCalled.Call()
				if exchange != r.Exchange {
					t.Fatalf("unexpected exchange: got %q want %q", exchange, r.Exchange)
				}
				if key != r.Key {
					t.Fatalf("unexpected key: got %q, want %q", key, r.Key)
				}
				testutils.Compare(t, "unexpected publishing", pbl, tc.expectedPbl)
				return nil
			}
			err := r.Retry(ctx, tc.dlv)
			if err == nil {
				t.Fatal("no error")
			}
			if !errors.IsIgnored(err) {
				t.Fatal("not ignored")
			}
			a := GetErrorAcknowledger(err)
			if a != Ack {
				t.Fatalf("unexpected acknowledger: got %v, want %v", a, Ack)
			}
			pCalled.AssertCalled(t)
		})
	}
}

func TestRetryerMaxReached(t *testing.T) {
	ctx := context.Background()
	dlv := amqp.Delivery{
		Headers: amqp.Table{
			retryHeaderAttempts: int64(10),
		},
		Body: []byte("test"),
	}
	r := &Retryer{
		Max: 10,
	}
	err := r.Retry(ctx, dlv)
	if err == nil {
		t.Fatal("no error")
	}
	a := GetErrorAcknowledger(err)
	if a != NackDiscard {
		t.Fatalf("unexpected acknowledger: got %v, want %v", a, NackDiscard)
	}
}

func TestRetryerMaxReachedAck(t *testing.T) {
	ctx := context.Background()
	dlv := amqp.Delivery{
		Headers: amqp.Table{
			retryHeaderAttempts: int64(10),
		},
		Body: []byte("test"),
	}
	r := &Retryer{
		Max:    10,
		MaxAck: true,
	}
	err := r.Retry(ctx, dlv)
	if err == nil {
		t.Fatal("no error")
	}
	a := GetErrorAcknowledger(err)
	if a != Ack {
		t.Fatalf("unexpected acknowledger: got %v, want %v", a, Ack)
	}
}

func TestRetryerInfinite(t *testing.T) {
	ctx := context.Background()
	dlv := amqp.Delivery{
		Headers: amqp.Table{
			retryHeaderAttempts: int64(10),
		},
		Body: []byte("test"),
	}
	r := &Retryer{
		Delay: 1 * time.Minute,
	}
	expectedPbl := amqp.Publishing{
		Body:       dlv.Body,
		Headers:    amqp.Table{},
		Expiration: "60000",
	}
	var pCalled testutils.CallCounter
	r.Producer = func(ctx context.Context, exchange, key string, mandatory, immediate bool, pbl amqp.Publishing) error {
		pCalled.Call()
		if exchange != r.Exchange {
			t.Fatalf("unexpected exchange: got %q want %q", exchange, r.Exchange)
		}
		if key != r.Key {
			t.Fatalf("unexpected key: got %q, want %q", key, r.Key)
		}
		testutils.Compare(t, "unexpected publishing", pbl, expectedPbl)
		return nil
	}
	err := r.Retry(ctx, dlv)
	if err == nil {
		t.Fatal("no error")
	}
	a := GetErrorAcknowledger(err)
	if a != Ack {
		t.Fatalf("unexpected acknowledger: got %v, want %v", a, Ack)
	}
	pCalled.AssertCalled(t)
}

func TestRetryerErrorProducer(t *testing.T) {
	ctx := context.Background()
	dlv := amqp.Delivery{
		Body: []byte("test"),
	}
	r := &Retryer{
		Max:      10,
		Delay:    1 * time.Minute,
		Exchange: "exchange",
		Key:      "key",
	}
	r.Producer = func(ctx context.Context, exchange, key string, mandatory, immediate bool, pbl amqp.Publishing) error {
		return errors.New("error")
	}
	err := r.Retry(ctx, dlv)
	if err == nil {
		t.Fatal("no error")
	}
}
