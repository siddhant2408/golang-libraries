package kafkautils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/siddhant2408/golang-libraries/timeutils"
)

func TestCopyMessage(t *testing.T) {
	msg1 := kafka.Message{
		Key:   []byte("test"),
		Value: []byte("test"),
		Headers: []kafka.Header{
			{
				Key:   "test",
				Value: []byte("test"),
			},
		},
		Time: timeutils.Now(),
	}
	msg2 := CopyMessage(msg1)
	testutils.Compare(t, "unexpected message", msg2, msg1)
}

func TestBytesTruncateConvertString(t *testing.T) {
	s := "test"
	v := bytesTruncateConvert([]byte(s))
	sRes, ok := v.(string)
	if !ok {
		t.Fatal("not a string")
	}
	if sRes != s {
		t.Fatal("not equal")
	}
}

func TestBytesTruncateConvertBytes(t *testing.T) {
	b := []byte{0xc3, 0x28}
	v := bytesTruncateConvert(b)
	bRes, ok := v.([]byte)
	if !ok {
		t.Fatal("not a []byte")
	}
	if !bytes.Equal(bRes, b) {
		t.Fatal("not equal")
	}
}

func TestBytesTruncateConvertMaxSize(t *testing.T) {
	b := []byte(strings.Repeat("a", 10000))
	v := bytesTruncateConvert(b)
	sRes, ok := v.(string)
	if !ok {
		t.Fatal("not a string")
	}
	expected := strings.Repeat("a", 512)
	if sRes != expected {
		t.Fatal("not equal")
	}
}
