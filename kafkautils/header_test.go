package kafkautils

import (
	"bytes"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestAddHeader(t *testing.T) {
	hs := []kafka.Header{
		{
			Key:   "a",
			Value: []byte("value a"),
		},
	}
	hs = AddHeader(hs, "b", []byte("value b"))
	expected := []kafka.Header{
		{
			Key:   "a",
			Value: []byte("value a"),
		},
		{
			Key:   "b",
			Value: []byte("value b"),
		},
	}
	testutils.Compare(t, "unexpected headers", hs, expected)
}

func TestSetHeaderNotExists(t *testing.T) {
	hs := []kafka.Header{
		{
			Key:   "a",
			Value: []byte("value a"),
		},
	}
	hs = SetHeader(hs, "b", []byte("value b"))
	expected := []kafka.Header{
		{
			Key:   "a",
			Value: []byte("value a"),
		},
		{
			Key:   "b",
			Value: []byte("value b"),
		},
	}
	testutils.Compare(t, "unexpected headers", hs, expected)
}

func TestSetHeaderAlreadyExists(t *testing.T) {
	hs := []kafka.Header{
		{
			Key:   "a",
			Value: []byte("value 1"),
		},
	}
	hs = SetHeader(hs, "a", []byte("value 2"))
	expected := []kafka.Header{
		{
			Key:   "a",
			Value: []byte("value 2"),
		},
	}
	testutils.Compare(t, "unexpected headers", hs, expected)
}

func TestGetHeaderFound(t *testing.T) {
	hs := []kafka.Header{
		{
			Key:   "a",
			Value: []byte("value a"),
		},
	}
	value, ok := GetHeader(hs, "a")
	if !ok {
		t.Fatal("not found")
	}
	expectedValue := []byte("value a")
	if !bytes.Equal(value, expectedValue) {
		t.Fatal("not equal")
	}
}

func TestGetHeaderNotFound(t *testing.T) {
	hs := []kafka.Header{
		{
			Key:   "a",
			Value: []byte("value a"),
		},
	}
	_, ok := GetHeader(hs, "b")
	if ok {
		t.Fatal("found")
	}
}

func TestGetHeaderAll(t *testing.T) {
	hs := []kafka.Header{
		{
			Key:   "a",
			Value: []byte("value 1"),
		},
		{
			Key:   "a",
			Value: []byte("value 2"),
		},
		{
			Key:   "b",
			Value: []byte("value 3"),
		},
	}
	values := GetHeaderAll(hs, "a")
	expectedValues := [][]byte{[]byte("value 1"), []byte("value 2")}
	testutils.Compare(t, "unexpected values", values, expectedValues)
}

func TestDeleteHeader(t *testing.T) {
	hs := []kafka.Header{
		{
			Key:   "a",
			Value: []byte("value 1"),
		},
		{
			Key:   "b",
			Value: []byte("value 2"),
		},
	}
	hs = DeleteHeader(hs, "a")
	expected := []kafka.Header{
		{
			Key:   "b",
			Value: []byte("value 2"),
		},
	}
	testutils.Compare(t, "unexpected headers", hs, expected)
}
