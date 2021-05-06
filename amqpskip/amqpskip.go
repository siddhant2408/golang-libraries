// Package amqpskip provides utilities to manage skipped organizations/applications in an AMQP consumer.
package amqpskip

import (
	"context"

	"github.com/siddhant2408/golang-libraries/amqputils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
	"github.com/streadway/amqp"
)

const headerOrganizationID = "organization-id"

// Checker checks the skipped organization.
type Checker struct {
	Applications []string
	IsSkipped    func(ctx context.Context, orgID int64, apps ...string) (skipped bool, found bool, err error)
	Forwarder    func(ctx context.Context, dlv amqp.Delivery, orgID int64) error
}

// Check checks if an organization is skipped.
// It requires to have the organization ID in the "organization-id" header with the type int64.
// If the organization is not skipped, it returns not error.
// If the organizations is skipped, it forwards the message, and return an error that is ignored and acknowledge the current message.
func (c *Checker) Check(ctx context.Context, dlv amqp.Delivery) (err error) {
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "amqpskip.checker", &err)
	defer spanFinish()
	orgID, err := c.getOrganizationID(dlv)
	if err != nil {
		err = amqputils.ErrorWithAcknowledger(err, amqputils.NackDiscard)
		err = errors.Wrap(err, "get organization ID")
		return err
	}
	span.SetTag("organization.id", orgID)
	skipped, err := c.isSkipped(ctx, orgID)
	if err != nil {
		return errors.Wrap(err, "is skipped")
	}
	span.SetTag("skipped", skipped)
	if !skipped {
		return nil
	}
	err = c.Forwarder(ctx, dlv, orgID)
	if err != nil {
		return errors.Wrap(err, "forward")
	}
	err = errors.New("skipped")
	err = errors.Ignore(err)
	err = amqputils.ErrorWithAcknowledger(err, amqputils.Ack)
	return err
}

func (c *Checker) getOrganizationID(dlv amqp.Delivery) (int64, error) {
	v, ok := dlv.Headers[headerOrganizationID]
	if !ok {
		return 0, errors.Newf("missing header: %q", headerOrganizationID)
	}
	// Only check int64, because AMQP supports strict typing.
	// If you're reading this comment and you have a problem, then the bug is in the application that is publishing the message.
	// Don't change this code.
	id, ok := v.(int64)
	if !ok {
		return 0, errors.Newf("unsupported header type: %T", v)
	}
	return id, nil
}

func (c *Checker) isSkipped(ctx context.Context, orgID int64) (bool, error) {
	skipped, found, err := c.IsSkipped(ctx, orgID, c.Applications...)
	if err != nil {
		return false, err
	}
	if !found {
		err := errors.New("organization not found")
		err = errors.WithTagInt64(err, "organization.id", orgID)
		err = amqputils.ErrorWithAcknowledger(err, amqputils.NackDiscard)
		return false, err
	}
	return skipped, nil
}

// Forwarder forwards messages.
type Forwarder struct {
	Topology    func(orgID int64) (tpi TopologyInit, exchange string, key string)
	ChannelPool func(context.Context, func(context.Context, *amqp.Channel) error) error
	Producer    amqputils.Producer
}

// Forward forwards the message for the given organization ID.
func (f *Forwarder) Forward(ctx context.Context, dlv amqp.Delivery, orgID int64) (err error) {
	_, spanFinish := tracingutils.StartChildSpan(&ctx, "amqpskip.forwarder", &err)
	defer spanFinish()
	tpi, exchange, key := f.Topology(orgID)
	err = f.ChannelPool(ctx, tpi)
	if err != nil {
		return errors.Wrap(err, "channel pool topology")
	}
	err = amqputils.Reproduce(ctx, f.Producer, exchange, key, false, false, dlv)
	if err != nil {
		return errors.Wrap(err, "reproduce")
	}
	return nil
}

// TopologyInit is a function type for *amqputils.Topology.Init().
type TopologyInit func(context.Context, *amqp.Channel) error

// ConsumerProcessor checks that the organization is not skipped and forwards to a ConsumerProcessor.
type ConsumerProcessor struct {
	amqputils.ConsumerProcessor
	Checker func(context.Context, amqp.Delivery) error
}

// Process implements ConsumerProcessor.
func (p *ConsumerProcessor) Process(ctx context.Context, dlv amqp.Delivery) (err error) {
	_, spanFinish := tracingutils.StartChildSpan(&ctx, "amqpskip.consumer_processor", &err)
	defer spanFinish()
	err = p.Checker(ctx, dlv)
	if err != nil {
		return errors.Wrap(err, "check")
	}
	return p.ConsumerProcessor(ctx, dlv)
}
