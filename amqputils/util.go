package amqputils

import (
	"unicode/utf8"

	"github.com/streadway/amqp"
)

const bodyTruncateSize = 512

func bodyTruncateConvert(body []byte) interface{} {
	isString := utf8.Valid(body)
	if len(body) > bodyTruncateSize {
		body = body[:bodyTruncateSize]
	}
	if isString {
		return string(body)
	}
	return body
}

func copyValue(o interface{}) interface{} {
	switch o := o.(type) {
	case amqp.Table:
		return copyTable(o)
	case []interface{}:
		return copySlice(o)
	case []byte:
		return copyBytes(o)
	}
	return o
}

func copyTable(o amqp.Table) amqp.Table {
	c := make(amqp.Table, len(o))
	for k, v := range o {
		c[k] = copyValue(v)
	}
	return c
}

func copySlice(o []interface{}) []interface{} {
	c := make([]interface{}, len(o))
	for i, v := range o {
		c[i] = copyValue(v)
	}
	return c
}

func copyBytes(o []byte) []byte {
	c := make([]byte, len(o))
	copy(c, o)
	return c
}

func deliveryToPublishing(dlv amqp.Delivery) amqp.Publishing {
	pbl := amqp.Publishing{
		Headers:         copyTable(dlv.Headers),
		ContentType:     dlv.ContentType,
		ContentEncoding: dlv.ContentEncoding,
		DeliveryMode:    dlv.DeliveryMode,
		Priority:        dlv.Priority,
		CorrelationId:   dlv.CorrelationId,
		ReplyTo:         dlv.ReplyTo,
		Expiration:      dlv.Expiration,
		MessageId:       dlv.MessageId,
		Timestamp:       dlv.Timestamp,
		Type:            dlv.Type,
		UserId:          dlv.UserId,
		AppId:           dlv.AppId,
		Body:            dlv.Body,
	}
	// It doesn't make sense to keep the "x-death" header.
	delete(pbl.Headers, "x-death")
	return pbl
}
