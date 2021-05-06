package amqputils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
	"github.com/streadway/amqp"
)

func TestBodyTruncateConvertString(t *testing.T) {
	s := "test"
	v := bodyTruncateConvert([]byte(s))
	sRes, ok := v.(string)
	if !ok {
		t.Fatal("not a string")
	}
	if sRes != s {
		t.Fatal("not equal")
	}
}

func TestBodyTruncateConvertBytes(t *testing.T) {
	b := []byte{0xc3, 0x28}
	v := bodyTruncateConvert(b)
	bRes, ok := v.([]byte)
	if !ok {
		t.Fatal("not a []byte")
	}
	if !bytes.Equal(bRes, b) {
		t.Fatal("not equal")
	}
}

func TestBodyTruncateConvertMaxSize(t *testing.T) {
	b := []byte(strings.Repeat("a", 10000))
	v := bodyTruncateConvert(b)
	sRes, ok := v.(string)
	if !ok {
		t.Fatal("not a string")
	}
	expected := strings.Repeat("a", 512)
	if sRes != expected {
		t.Fatal("not equal")
	}
}

func TestCopyValue(t *testing.T) {
	v := amqp.Table{
		"string": "test",
		"slice": []interface{}{
			"test",
		},
		"bytes": []byte("test"),
	}
	res := copyValue(v)
	testutils.Compare(t, "unexpected copy", res, v)
}

func TestDeliveryToPublishing(t *testing.T) {
	dlv := amqp.Delivery{
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
	}
	pbl := deliveryToPublishing(dlv)
	expected := amqp.Publishing{
		Headers: amqp.Table{
			"foo": "bar",
		},
		Body: []byte("test"),
	}
	testutils.Compare(t, "unexpected publishing", pbl, expected)
}
