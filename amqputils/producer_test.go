package amqputils_test

import (
	"context"
	"testing"
	"time"

	"github.com/siddhant2408/golang-libraries/amqptest"
	"github.com/siddhant2408/golang-libraries/amqputils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/streadway/amqp"
)

func TestSimpleProducer(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	p := &amqputils.SimpleProducer{
		Channel: amqputils.NewChannelGetterConnection(conn),
	}
	defer p.Close() //nolint:errcheck
	pbl := amqp.Publishing{
		Headers: amqp.Table{
			"foo": "bar",
		},
	}
	err := p.Produce(ctx, "amq.direct", "test", false, false, pbl)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = p.Produce(ctx, "amq.direct", "test", false, false, pbl)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestSimpleProducerConfirm(t *testing.T) {
	ctx := context.Background()
	conn := amqptest.NewConnection(t, testVhost)
	p := &amqputils.SimpleProducer{
		Channel: amqputils.NewChannelGetterConnection(conn),
		Confirm: true,
	}
	defer p.Close() //nolint:errcheck
	pbl := amqp.Publishing{}
	err := p.Produce(ctx, "", "", false, false, pbl)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestSimpleProducerErrorLock(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	p := &amqputils.SimpleProducer{}
	defer p.Close() //nolint:errcheck
	err := p.Produce(ctx, "", "_test", false, false, amqp.Publishing{
		Body: []byte("test"),
	})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSimpleProducerErrorOpenChannel(t *testing.T) {
	p := &amqputils.SimpleProducer{
		Channel: func(context.Context) (*amqp.Channel, error) {
			return nil, errors.New("error")
		},
	}
	defer p.Close() //nolint:errcheck
	err := p.Produce(context.Background(), "", "_test", false, false, amqp.Publishing{
		Body: []byte("test"),
	})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestTimeoutProducer(t *testing.T) {
	ctx := context.Background()
	p := &amqputils.TimeoutProducer{
		Producer: func(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			return nil
		},
		Timeout: 1 * time.Hour,
	}
	err := p.Produce(ctx, "", "", false, false, amqp.Publishing{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestTimeoutProducerErrorTimeout(t *testing.T) {
	ctx := context.Background()
	p := &amqputils.TimeoutProducer{
		Producer: func(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			<-ctx.Done()
			return ctx.Err()
		},
		Timeout: 1 * time.Microsecond,
	}
	err := p.Produce(ctx, "", "", false, false, amqp.Publishing{})
	if err == nil {
		t.Fatal("no error")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal("not deadline exceeded")
	}
}

func TestMultiProducer(t *testing.T) {
	ctx := context.Background()
	var pCalled testutils.CallCounter
	p := func(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
		pCalled.Call()
		return nil
	}
	mp := amqputils.NewMultiProducer([]amqputils.Producer{p})
	pbl := amqp.Publishing{}
	err := mp.Produce(ctx, "", "", false, false, pbl)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	pCalled.AssertCalled(t)
}

func TestMultiProducerError(t *testing.T) {
	ctx := context.Background()
	p := func(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
		return errors.New("error")
	}
	mp := amqputils.NewMultiProducer([]amqputils.Producer{p})
	pbl := amqp.Publishing{}
	err := mp.Produce(ctx, "", "", false, false, pbl)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestMultiProducerContextDone(t *testing.T) {
	p := amqputils.NewMultiProducer(nil)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	err := p.Produce(ctx, "", "_test", false, false, amqp.Publishing{
		Body: []byte("test"),
	})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestBufferedProducer(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	msg := amqp.Publishing{
		Body: []byte(`test`),
	}
	var tpCalled testutils.CallCounter
	tp := func(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
		tpCalled.Call()
		cancel()
		return nil
	}
	p := amqputils.NewBufferProducer(tp, 1, nil)
	err := p.Produce(ctx, "test", "test", false, false, msg)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	p.Run(ctx)
	tpCalled.AssertCalled(t)
}

func TestBufferedProducerDrain(t *testing.T) {
	ctx := context.Background()
	msg := amqp.Publishing{
		Body: []byte(`test`),
	}
	var tpCalled testutils.CallCounter
	tp := func(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
		tpCalled.Call()
		return nil
	}
	p := amqputils.NewBufferProducer(tp, 1, nil)
	err := p.Produce(ctx, "test", "test", false, false, msg)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	p.Drain(ctx)
	tpCalled.AssertCalled(t)
}

func TestBufferedProducerErrorProducer(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	msg := amqp.Publishing{
		Body: []byte(`test`),
	}
	tp := func(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
		return errors.New("error")
	}
	var errFuncCalled testutils.CallCounter
	errFunc := func(ctx context.Context, err error) {
		errFuncCalled.Call()
		cancel()
	}
	p := amqputils.NewBufferProducer(tp, 1, errFunc)
	err := p.Produce(ctx, "test", "test", false, false, msg)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	p.Run(ctx)
	errFuncCalled.AssertCalled(t)
}

func TestBufferedProducerErrorFull(t *testing.T) {
	ctx := context.Background()
	msg := amqp.Publishing{
		Body: []byte(`test`),
	}
	p := amqputils.NewBufferProducer(nil, 0, nil)
	err := p.Produce(ctx, "test", "test", false, false, msg)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestErrorProducer(t *testing.T) {
	ctx := context.Background()
	msg := amqp.Publishing{
		Body: []byte(`test`),
	}
	var tpCalled testutils.CallCounter
	tp := func(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
		tpCalled.Call()
		return nil
	}
	p := amqputils.NewErrorProducer(tp, nil)
	p.Produce(ctx, "test", "test", false, false, msg)
	tpCalled.AssertCalled(t)
}

func TestErrorProducerError(t *testing.T) {
	ctx := context.Background()
	msg := amqp.Publishing{
		Body: []byte(`test`),
	}
	tp := func(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
		return errors.New("error")
	}
	var errFuncCalled testutils.CallCounter
	errFunc := func(ctx context.Context, err error) {
		errFuncCalled.Call()
	}
	p := amqputils.NewErrorProducer(tp, errFunc)
	p.Produce(ctx, "test", "test", false, false, msg)
	errFuncCalled.AssertCalled(t)
}

func TestMultiConfirmProducer(t *testing.T) {
	conn := amqptest.NewConnection(t, testVhost)
	cg := amqputils.NewChannelGetterConnection(conn)
	p, closeP := amqputils.NewMultiConfirmProducer(cg, 10)
	if p == nil {
		t.Fatal("nil")
	}
	closeP(func(err error) {
		testutils.FatalErr(t, err)
	})
}

func TestReproduce(t *testing.T) {
	expectedMsg := amqp.Publishing{
		Headers: amqp.Table{
			"foo": "bar",
		},
		Body: []byte("test"),
	}
	var pCalled testutils.CallCounter
	p := func(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
		pCalled.Call()
		testutils.Compare(t, "unexpected message", msg, expectedMsg)
		return nil
	}
	err := amqputils.Reproduce(context.Background(), p, "", "_test", false, false, amqp.Delivery{
		Headers: amqp.Table{
			"foo": "bar",
			"x-death": []interface{}{
				amqp.Table{
					"exchange": "a",
					"routing-keys": []interface{}{
						"b",
					},
				},
			},
		},
		Body: []byte("test"),
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	pCalled.AssertCalled(t)
}
