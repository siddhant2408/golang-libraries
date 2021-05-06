package kafkautils

import (
	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/errors"
)

func wrapErrorValue(err error, key string, val interface{}) error {
	return errors.WithValue(err, "kafka."+key, val)
}

func wrapErrorValueMessage(err error, msg kafka.Message) error {
	if !msg.Time.IsZero() {
		err = wrapErrorValue(err, "message.time", msg.Time)
	}
	if len(msg.Value) > 0 {
		err = wrapErrorValue(err, "message.value", bytesTruncateConvert(msg.Value))
	}
	if len(msg.Key) > 0 {
		err = wrapErrorValue(err, "message.key", bytesTruncateConvert(msg.Key))
	}
	if len(msg.Topic) > 0 {
		err = wrapErrorValue(err, "message.topic", msg.Topic)
	}
	// TODO headers
	return err
}

func wrapErrorValueMessageConsumer(err error, msg kafka.Message) error {
	err = wrapErrorValue(err, "message.offset", msg.Offset)
	err = wrapErrorValue(err, "message.partition", msg.Partition)
	err = wrapErrorValue(err, "message.topic", msg.Topic)
	err = wrapErrorValueMessage(err, msg)
	return err
}
