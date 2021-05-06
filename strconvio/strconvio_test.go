package strconvio_test

import (
	"bytes"
	"io"
	"testing"

	. "github.com/siddhant2408/golang-libraries/strconvio"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestWriteBool(t *testing.T) {
	buf := new(bytes.Buffer)
	_, err := WriteBool(buf, true)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	s := buf.String()
	expected := "true"
	if s != expected {
		t.Fatalf("unexpected result: got %q, want %q", s, expected)
	}
}

func BenchmarkWriteBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = WriteBool(io.Discard, true)
	}
}

func TestWriteFloat(t *testing.T) {
	buf := new(bytes.Buffer)
	_, err := WriteFloat(buf, 123.456, 'f', -1, 64)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	s := buf.String()
	expected := "123.456"
	if s != expected {
		t.Fatalf("unexpected result: got %q, want %q", s, expected)
	}
}

func BenchmarkWriteFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = WriteFloat(io.Discard, 123.456, 'f', -1, 64)
	}
}

func TestWriteInt(t *testing.T) {
	buf := new(bytes.Buffer)
	_, err := WriteInt(buf, 123, 10)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	s := buf.String()
	expected := "123"
	if s != expected {
		t.Fatalf("unexpected result: got %q, want %q", s, expected)
	}
}

func BenchmarkWriteInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = WriteInt(io.Discard, 123, 10)
	}
}

func TestWriteUint(t *testing.T) {
	buf := new(bytes.Buffer)
	_, err := WriteUint(buf, 123, 10)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	s := buf.String()
	expected := "123"
	if s != expected {
		t.Fatalf("unexpected result: got %q, want %q", s, expected)
	}
}

func BenchmarkWriteUint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = WriteUint(io.Discard, 123, 10)
	}
}
