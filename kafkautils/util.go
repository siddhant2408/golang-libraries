package kafkautils

import (
	"unicode/utf8"

	"github.com/segmentio/kafka-go"
)

// CopyMessage copies a message.
// It doesn't contain the fields that shouldn't be used in a Producer.
func CopyMessage(msg kafka.Message) kafka.Message {
	var copyMsg kafka.Message
	if msg.Key != nil {
		copyMsg.Key = make([]byte, len(msg.Key))
		copy(copyMsg.Key, msg.Key)
	}
	if msg.Value != nil {
		copyMsg.Value = make([]byte, len(msg.Value))
		copy(copyMsg.Value, msg.Value)
	}
	if msg.Headers != nil {
		copyMsg.Headers = make([]kafka.Header, len(msg.Headers))
		copy(copyMsg.Headers, msg.Headers)
	}
	copyMsg.Time = msg.Time
	return copyMsg
}

const bytesTruncateSize = 512

func bytesTruncateConvert(value []byte) interface{} {
	isString := utf8.Valid(value)
	if len(value) > bytesTruncateSize {
		value = value[:bytesTruncateSize]
	}
	if isString {
		return string(value)
	}
	return value
}
