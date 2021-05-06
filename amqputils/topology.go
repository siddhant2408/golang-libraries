package amqputils

import (
	"context"

	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
	"github.com/streadway/amqp"
)

// InitTopology initializes a topology.
func InitTopology(ctx context.Context, cg ChannelGetter, tp Topology) error {
	chn, err := cg(ctx)
	if err != nil {
		return errors.Wrap(err, "get channel")
	}
	defer chn.Close() //nolint:errcheck
	return tp.Init(ctx, chn)
}

// Topology represents an exchanges+queues topology.
type Topology struct {
	Exchanges []ExchangeConfig
	Queues    []QueueConfig
}

// Init initializes the topology.
func (tp Topology) Init(ctx context.Context, chn *amqp.Channel) error {
	for _, ec := range tp.Exchanges {
		err := ec.init(ctx, chn)
		if err != nil {
			err = wrapErrorValue(err, "exchange", ec.Name)
			err = errors.Wrap(err, "exchange")
			return err
		}
	}
	for _, qc := range tp.Queues {
		err := qc.init(ctx, chn)
		if err != nil {
			err = wrapErrorValue(err, "queue", qc.Name)
			err = errors.Wrap(err, "queue")
			return err
		}
	}
	return nil
}

// ExchangeConfig represents an exchange config.
type ExchangeConfig struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	Arguments  amqp.Table
	Bindings   []ExchangeBinding
}

func (ec ExchangeConfig) init(ctx context.Context, chn *amqp.Channel) error {
	err := ec.declare(ctx, chn)
	if err != nil {
		return errors.Wrap(err, "declare")
	}
	for _, eb := range ec.Bindings {
		err := eb.bind(ctx, chn, ec.Name)
		if err != nil {
			return errors.Wrap(err, "bind")
		}
	}
	return nil
}

func (ec ExchangeConfig) declare(ctx context.Context, chn *amqp.Channel) (err error) {
	span, spanFinish := startTraceChildSpan(&ctx, "exchange_declare", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, tracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.AppTypeRPC)
	opentracing_ext.SpanKindRPCClient.Set(span)
	setTraceSpanTag(span, "exchange", ec.Name)
	err = chn.ExchangeDeclare(ec.Name, ec.Type, ec.Durable, ec.AutoDelete, ec.Internal, false, ec.Arguments)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

// ExchangeBinding represents an exchange binding.
type ExchangeBinding struct {
	RoutingKey string
	Source     string
	Arguments  amqp.Table
}

func (eb ExchangeBinding) bind(ctx context.Context, chn *amqp.Channel, name string) (err error) {
	span, spanFinish := startTraceChildSpan(&ctx, "exchange_bind", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, tracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.AppTypeRPC)
	opentracing_ext.SpanKindRPCClient.Set(span)
	setTraceSpanTag(span, "exchange", name)
	setTraceSpanTag(span, "source", eb.Source)
	if eb.RoutingKey != "" {
		setTraceSpanTag(span, "routing_key", eb.RoutingKey)
	}
	err = chn.ExchangeBind(name, eb.RoutingKey, eb.Source, false, eb.Arguments)
	if err != nil {
		err = wrapErrorValue(err, "source", eb.Source)
		if eb.RoutingKey != "" {
			err = wrapErrorValue(err, "routing_key", eb.RoutingKey)
		}
		err = errors.Wrap(err, "")
		return err
	}
	return nil
}

// QueueConfig represents a queue config.
type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	Arguments  amqp.Table
	Bindings   []QueueBinding
}

func (qc QueueConfig) init(ctx context.Context, chn *amqp.Channel) error {
	err := qc.declare(ctx, chn)
	if err != nil {
		return errors.Wrap(err, "declare")
	}
	for _, qb := range qc.Bindings {
		err := qb.bind(ctx, chn, qc.Name)
		if err != nil {
			return errors.Wrap(err, "bind")
		}
	}
	return nil
}

func (qc QueueConfig) declare(ctx context.Context, chn *amqp.Channel) (err error) {
	span, spanFinish := startTraceChildSpan(&ctx, "queue_declare", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, tracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.AppTypeRPC)
	opentracing_ext.SpanKindRPCClient.Set(span)
	setTraceSpanTag(span, "queue", qc.Name)
	_, err = chn.QueueDeclare(qc.Name, qc.Durable, qc.AutoDelete, qc.Exclusive, false, qc.Arguments)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

// QueueBinding represents a queue binding.
type QueueBinding struct {
	RoutingKey string
	Exchange   string
	Arguments  amqp.Table
}

func (qb QueueBinding) bind(ctx context.Context, chn *amqp.Channel, name string) (err error) {
	span, spanFinish := startTraceChildSpan(&ctx, "queue_bind", &err)
	defer spanFinish()
	tracingutils.SetSpanServiceName(span, tracingExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.AppTypeRPC)
	opentracing_ext.SpanKindRPCClient.Set(span)
	setTraceSpanTag(span, "queue", name)
	setTraceSpanTag(span, "exchange", qb.Exchange)
	if qb.RoutingKey != "" {
		setTraceSpanTag(span, "routing_key", qb.RoutingKey)
	}
	err = chn.QueueBind(name, qb.RoutingKey, qb.Exchange, false, qb.Arguments)
	if err != nil {
		err = wrapErrorValue(err, "exchange", qb.Exchange)
		if qb.RoutingKey != "" {
			err = wrapErrorValue(err, "routing_key", qb.RoutingKey)
		}
		err = errors.Wrap(err, "")
		return err
	}
	return nil
}
