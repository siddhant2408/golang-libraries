package jsontracing

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestMarshalUnmarshal(t *testing.T) {
	ctx := context.Background()
	v1 := map[string]interface{}{
		"foo": "bar",
	}
	b, err := Marshal(ctx, v1)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var v2 interface{}
	err = Unmarshal(ctx, b, &v2)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	testutils.Compare(t, "unexpected value", v2, v1)
}

func TestMarshalError(t *testing.T) {
	ctx := context.Background()
	v := func() {}
	_, err := Marshal(ctx, v)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestUnmarshalError(t *testing.T) {
	ctx := context.Background()
	b := []byte("invalid")
	var v interface{}
	err := Unmarshal(ctx, b, &v)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestEncoderDecoder(t *testing.T) {
	ctx := context.Background()
	v1 := map[string]interface{}{
		"foo": "bar",
	}
	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	err := enc.Encode(ctx, v1)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var v2 interface{}
	dec := NewDecoder(buf)
	err = dec.Decode(ctx, &v2)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	testutils.Compare(t, "unexpected value", v2, v1)
}

func TestEncoderError(t *testing.T) {
	ctx := context.Background()
	v := func() {}
	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	err := enc.Encode(ctx, v)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestDecoderError(t *testing.T) {
	ctx := context.Background()
	r := strings.NewReader("invalid")
	var v interface{}
	dec := NewDecoder(r)
	err := dec.Decode(ctx, &v)
	if err == nil {
		t.Fatal("no error")
	}
}
