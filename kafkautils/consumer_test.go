package kafkautils

import (
	"context"
	"fmt"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestConsumer(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	fetch := func(context.Context) (kafka.Message, error) {
		return kafka.Message{}, nil
	}
	pr := func(ctx context.Context, msg kafka.Message) error {
		return nil
	}
	commit := func(context.Context, ...kafka.Message) error {
		cancel()
		return nil
	}
	r := &testFetchCommitter{
		fetch:  fetch,
		commit: commit,
	}
	c := &Consumer{
		Processor: pr,
	}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestConsumerFetchContextDone(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	fetch := func(ctx context.Context) (kafka.Message, error) {
		cancel()
		return kafka.Message{}, ctx.Err()
	}
	r := &testFetchCommitter{
		fetch: fetch,
	}
	c := &Consumer{}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestConsumerErrorFetch(t *testing.T) {
	ctx := context.Background()
	fetch := func(context.Context) (kafka.Message, error) {
		return kafka.Message{}, errors.New("error")
	}
	r := &testFetchCommitter{
		fetch: fetch,
	}
	c := &Consumer{}
	err := c.Consume(ctx, r)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestConsumerErrorCommit(t *testing.T) {
	ctx := context.Background()
	fetch := func(context.Context) (kafka.Message, error) {
		return kafka.Message{}, nil
	}
	pr := func(ctx context.Context, msg kafka.Message) error {
		return nil
	}
	commit := func(context.Context, ...kafka.Message) error {
		return errors.New("error")
	}
	r := &testFetchCommitter{
		fetch:  fetch,
		commit: commit,
	}
	c := &Consumer{
		Processor: pr,
	}
	err := c.Consume(ctx, r)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestConsumerErrorProcessRetry(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	fetch := func(context.Context) (kafka.Message, error) {
		return kafka.Message{}, nil
	}
	pr := func(ctx context.Context, msg kafka.Message) error {
		return errors.New("error")
	}
	var errFuncCalled testutils.CallCounter
	errFunc := func(ctx context.Context, err error) {
		errFuncCalled.Call()
		_ = err.Error()
	}
	retry := func(context.Context, ...kafka.Message) error {
		cancel()
		return nil
	}
	commit := func(context.Context, ...kafka.Message) error {
		return nil
	}
	r := &testFetchCommitter{
		fetch:  fetch,
		commit: commit,
	}
	c := &Consumer{
		Processor: pr,
		Retry:     retry,
		Error:     errFunc,
	}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	errFuncCalled.AssertCalled(t)
}

func TestConsumerErrorProcessRetryError(t *testing.T) {
	ctx := context.Background()
	fetch := func(context.Context) (kafka.Message, error) {
		return kafka.Message{}, nil
	}
	pr := func(ctx context.Context, msg kafka.Message) error {
		err := errors.New("error")
		err = ConsumerErrorWithHandler(err, ConsumerRetry)
		return err
	}
	var errFuncCalled testutils.CallCounter
	errFunc := func(ctx context.Context, err error) {
		errFuncCalled.Call()
		_ = err.Error()
	}
	retry := func(context.Context, ...kafka.Message) error {
		return errors.New("error")
	}
	commit := func(context.Context, ...kafka.Message) error {
		return nil
	}
	r := &testFetchCommitter{
		fetch:  fetch,
		commit: commit,
	}
	c := &Consumer{
		Processor: pr,
		Retry:     retry,
		Error:     errFunc,
	}
	err := c.Consume(ctx, r)
	if err == nil {
		t.Fatal("no error")
	}
	errFuncCalled.AssertCalled(t)
}

func TestConsumerErrorProcessDiscard(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	fetch := func(context.Context) (kafka.Message, error) {
		return kafka.Message{}, nil
	}
	pr := func(ctx context.Context, msg kafka.Message) error {
		err := errors.New("error")
		err = ConsumerErrorWithHandler(err, ConsumerDiscard)
		return err
	}
	var errFuncCalled testutils.CallCounter
	errFunc := func(ctx context.Context, err error) {
		errFuncCalled.Call()
		_ = err.Error()
	}
	discard := func(context.Context, ...kafka.Message) error {
		cancel()
		return nil
	}
	commit := func(context.Context, ...kafka.Message) error {
		return nil
	}
	r := &testFetchCommitter{
		fetch:  fetch,
		commit: commit,
	}
	c := &Consumer{
		Processor: pr,
		Discard:   discard,
		Error:     errFunc,
	}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	errFuncCalled.AssertCalled(t)
}

func TestConsumerErrorProcessDiscardNotDefined(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	fetch := func(context.Context) (kafka.Message, error) {
		return kafka.Message{}, nil
	}
	pr := func(ctx context.Context, msg kafka.Message) error {
		cancel()
		err := errors.New("error")
		err = ConsumerErrorWithHandler(err, ConsumerDiscard)
		return err
	}
	var errFuncCalled testutils.CallCounter
	errFunc := func(ctx context.Context, err error) {
		errFuncCalled.Call()
		_ = err.Error()
	}
	commit := func(context.Context, ...kafka.Message) error {
		return nil
	}
	r := &testFetchCommitter{
		fetch:  fetch,
		commit: commit,
	}
	c := &Consumer{
		Processor: pr,
		Error:     errFunc,
	}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	errFuncCalled.AssertCalled(t)
}

func TestConsumerErrorProcessDiscardError(t *testing.T) {
	ctx := context.Background()
	fetch := func(context.Context) (kafka.Message, error) {
		return kafka.Message{}, nil
	}
	pr := func(ctx context.Context, msg kafka.Message) error {
		err := errors.New("error")
		err = ConsumerErrorWithHandler(err, ConsumerDiscard)
		return err
	}
	var errFuncCalled testutils.CallCounter
	errFunc := func(ctx context.Context, err error) {
		errFuncCalled.Call()
		_ = err.Error()
	}
	discard := func(context.Context, ...kafka.Message) error {
		return errors.New("error")
	}
	commit := func(context.Context, ...kafka.Message) error {
		return nil
	}
	r := &testFetchCommitter{
		fetch:  fetch,
		commit: commit,
	}
	c := &Consumer{
		Processor: pr,
		Discard:   discard,
		Error:     errFunc,
	}
	err := c.Consume(ctx, r)
	if err == nil {
		t.Fatal("no error")
	}
	errFuncCalled.AssertCalled(t)
}

func TestConsumerErrorProcessNotTemporary(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	fetch := func(context.Context) (kafka.Message, error) {
		return kafka.Message{}, nil
	}
	pr := func(ctx context.Context, msg kafka.Message) error {
		err := errors.New("error")
		err = errors.WithTemporary(err, false)
		return err
	}
	var errFuncCalled testutils.CallCounter
	errFunc := func(ctx context.Context, err error) {
		errFuncCalled.Call()
		_ = err.Error()
	}
	discard := func(context.Context, ...kafka.Message) error {
		cancel()
		return nil
	}
	commit := func(context.Context, ...kafka.Message) error {
		return nil
	}
	r := &testFetchCommitter{
		fetch:  fetch,
		commit: commit,
	}
	c := &Consumer{
		Processor: pr,
		Discard:   discard,
		Error:     errFunc,
	}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	errFuncCalled.AssertCalled(t)
}

func TestConsumerErrorProcessNoop(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	fetch := func(context.Context) (kafka.Message, error) {
		return kafka.Message{}, nil
	}
	pr := func(ctx context.Context, msg kafka.Message) error {
		cancel()
		err := errors.New("error")
		err = ConsumerErrorWithHandler(err, ConsumerNoop)
		return err
	}
	var errFuncCalled testutils.CallCounter
	errFunc := func(ctx context.Context, err error) {
		errFuncCalled.Call()
		_ = err.Error()
	}
	commit := func(context.Context, ...kafka.Message) error {
		return nil
	}
	r := &testFetchCommitter{
		fetch:  fetch,
		commit: commit,
	}
	c := &Consumer{
		Processor: pr,
		Error:     errFunc,
	}
	err := c.Consume(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	errFuncCalled.AssertCalled(t)
}

func TestConsumerErrorWithHandlerNil(t *testing.T) {
	err := ConsumerErrorWithHandler(nil, ConsumerNoop)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestConsumerHandlerErrorFormat(t *testing.T) {
	err := errors.New("error")
	err = ConsumerErrorWithHandler(err, ConsumerNoop)
	_ = err.Error()
	_ = fmt.Sprint(err)
}
