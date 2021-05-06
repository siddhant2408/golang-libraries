package amqpskip

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/amqputils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/streadway/amqp"
)

func TestCheckerSkipped(t *testing.T) {
	ctx := context.Background()
	orgID := int64(123)
	dlv := amqp.Delivery{
		Headers: amqp.Table{
			headerOrganizationID: orgID,
		},
	}
	apps := []string{"test"}
	var isSkippedCalled testutils.CallCounter
	isSkipped := func(ctx context.Context, orgIDSkipped int64, appsSKipped ...string) (skipped bool, found bool, err error) {
		isSkippedCalled.Call()
		if orgIDSkipped != orgID {
			t.Fatalf("unexpected organization ID: got %d, want %d", orgIDSkipped, orgID)
		}
		testutils.Compare(t, "unexpected applications", appsSKipped, apps)
		return true, true, nil
	}
	var forwarderCalled testutils.CallCounter
	forwarder := func(ctx context.Context, dlv amqp.Delivery, orgIDForward int64) error {
		forwarderCalled.Call()
		if orgIDForward != orgID {
			t.Fatalf("unexpected organization ID: got %d, want %d", orgIDForward, orgID)
		}
		return nil
	}
	c := &Checker{
		Applications: apps,
		IsSkipped:    isSkipped,
		Forwarder:    forwarder,
	}
	err := c.Check(ctx, dlv)
	if err == nil {
		t.Fatal("no error")
	}
	ignore := errors.IsIgnored(err)
	if !ignore {
		t.Fatal("not ignored")
	}
	ack := amqputils.GetErrorAcknowledger(err)
	if ack != amqputils.Ack {
		t.Fatalf("unexpected acknowledger: got %v, want %v", ack, amqputils.Ack)
	}
	isSkippedCalled.AssertCalled(t)
	forwarderCalled.AssertCalled(t)
}

func TestCheckerNotSkipped(t *testing.T) {
	ctx := context.Background()
	orgID := int64(123)
	dlv := amqp.Delivery{
		Headers: amqp.Table{
			headerOrganizationID: orgID,
		},
	}
	apps := []string{"test"}
	var isSkippedCalled testutils.CallCounter
	isSkipped := func(ctx context.Context, orgIDSkipped int64, appsSKipped ...string) (skipped bool, found bool, err error) {
		isSkippedCalled.Call()
		if orgIDSkipped != orgID {
			t.Fatalf("unexpected organization ID: got %d, want %d", orgIDSkipped, orgID)
		}
		testutils.Compare(t, "unexpected applications", appsSKipped, apps)
		return false, true, nil
	}
	c := &Checker{
		Applications: apps,
		IsSkipped:    isSkipped,
	}
	err := c.Check(ctx, dlv)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	isSkippedCalled.AssertCalled(t)
}

func TestCheckerErrorHeaderMissing(t *testing.T) {
	ctx := context.Background()
	dlv := amqp.Delivery{}
	c := &Checker{}
	err := c.Check(ctx, dlv)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCheckerHeaderWrongType(t *testing.T) {
	ctx := context.Background()
	dlv := amqp.Delivery{
		Headers: amqp.Table{
			headerOrganizationID: "invalid",
		},
	}
	c := &Checker{}
	err := c.Check(ctx, dlv)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCheckerErrorIsSkippedNotFound(t *testing.T) {
	ctx := context.Background()
	orgID := int64(123)
	dlv := amqp.Delivery{
		Headers: amqp.Table{
			headerOrganizationID: orgID,
		},
	}
	isSkipped := func(ctx context.Context, orgIDSkipped int64, appsSKipped ...string) (skipped bool, found bool, err error) {
		return false, false, nil
	}
	c := &Checker{
		IsSkipped: isSkipped,
	}
	err := c.Check(ctx, dlv)
	if err == nil {
		t.Fatal("no error")
	}
	ack := amqputils.GetErrorAcknowledger(err)
	if ack != amqputils.NackDiscard {
		t.Fatalf("unexpected acknowledger: got %v, want %v", ack, amqputils.NackDiscard)
	}
}

func TestCheckerErrorIsSkipped(t *testing.T) {
	ctx := context.Background()
	orgID := int64(123)
	dlv := amqp.Delivery{
		Headers: amqp.Table{
			headerOrganizationID: orgID,
		},
	}
	isSkipped := func(ctx context.Context, orgIDSkipped int64, appsSKipped ...string) (skipped bool, found bool, err error) {
		return false, false, errors.New("error")
	}
	c := &Checker{
		IsSkipped: isSkipped,
	}
	err := c.Check(ctx, dlv)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCheckerErrorForwarder(t *testing.T) {
	ctx := context.Background()
	orgID := int64(123)
	dlv := amqp.Delivery{
		Headers: amqp.Table{
			headerOrganizationID: orgID,
		},
	}
	apps := []string{"test"}
	isSkipped := func(ctx context.Context, orgIDSkipped int64, appsSKipped ...string) (skipped bool, found bool, err error) {
		return true, true, nil
	}
	forwarder := func(ctx context.Context, dlv amqp.Delivery, orgIDForward int64) error {
		return errors.New("error")
	}
	c := &Checker{
		Applications: apps,
		IsSkipped:    isSkipped,
		Forwarder:    forwarder,
	}
	err := c.Check(ctx, dlv)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestFowarder(t *testing.T) {
	ctx := context.Background()
	dlv := amqp.Delivery{}
	orgID := int64(123)
	var tpiCalled testutils.CallCounter
	tpi := func(context.Context, *amqp.Channel) error { //nolint:unparam // The error is always nil.
		tpiCalled.Call()
		return nil
	}
	var topologyCalled testutils.CallCounter
	topology := func(orgID int64) (_ TopologyInit, exchange string, key string) {
		topologyCalled.Call()
		return tpi, "exchange", "key"
	}
	channelPool := func(ctx context.Context, f func(context.Context, *amqp.Channel) error) error {
		chn := new(amqp.Channel)
		return f(ctx, chn)
	}
	var producerCalled testutils.CallCounter
	producer := func(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
		producerCalled.Call()
		return nil
	}
	f := &Forwarder{
		Topology:    topology,
		ChannelPool: channelPool,
		Producer:    producer,
	}
	err := f.Forward(ctx, dlv, orgID)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	topologyCalled.AssertCalled(t)
	tpiCalled.AssertCalled(t)
	producerCalled.AssertCalled(t)
}

func TestFowarderErrorTopologyInit(t *testing.T) {
	ctx := context.Background()
	dlv := amqp.Delivery{}
	orgID := int64(123)
	tpi := func(context.Context, *amqp.Channel) error {
		return errors.New("error")
	}
	topology := func(orgID int64) (_ TopologyInit, exchange string, key string) {
		return tpi, "exchange", "key"
	}
	channelPool := func(ctx context.Context, f func(context.Context, *amqp.Channel) error) error {
		chn := new(amqp.Channel)
		return f(ctx, chn)
	}
	f := &Forwarder{
		Topology:    topology,
		ChannelPool: channelPool,
	}
	err := f.Forward(ctx, dlv, orgID)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestFowarderErrorProducer(t *testing.T) {
	ctx := context.Background()
	dlv := amqp.Delivery{}
	orgID := int64(123)
	tpi := func(context.Context, *amqp.Channel) error {
		return nil
	}
	topology := func(orgID int64) (_ TopologyInit, exchange string, key string) {
		return tpi, "exchange", "key"
	}
	channelPool := func(ctx context.Context, f func(context.Context, *amqp.Channel) error) error {
		chn := new(amqp.Channel)
		return f(ctx, chn)
	}
	producer := func(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
		return errors.New("error")
	}
	f := &Forwarder{
		Topology:    topology,
		ChannelPool: channelPool,
		Producer:    producer,
	}
	err := f.Forward(ctx, dlv, orgID)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestConsumerProcessor(t *testing.T) {
	ctx := context.Background()
	dlv := amqp.Delivery{}
	c := func(ctx context.Context, dlv amqp.Delivery) error {
		return nil
	}
	var tpCalled testutils.CallCounter
	tp := func(ctx context.Context, dlv amqp.Delivery) error {
		tpCalled.Call()
		return nil
	}
	p := &ConsumerProcessor{
		ConsumerProcessor: tp,
		Checker:           c,
	}
	err := p.Process(ctx, dlv)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	tpCalled.AssertCalled(t)
}

func TestConsumerProcessorErrorCheck(t *testing.T) {
	ctx := context.Background()
	dlv := amqp.Delivery{}
	c := func(ctx context.Context, dlv amqp.Delivery) error {
		return errors.New("error")
	}
	p := &ConsumerProcessor{
		Checker: c,
	}
	err := p.Process(ctx, dlv)
	if err == nil {
		t.Fatal("no error")
	}
}
