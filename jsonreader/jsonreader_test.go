package jsonreader

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestRead(t *testing.T) {
	w := New("test", func(enc *json.Encoder) {
		enc.SetEscapeHTML(false)
	})
	b, err := io.ReadAll(w)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	expected := []byte("\"test\"\n")
	if !bytes.Equal(b, expected) {
		t.Fatalf("unexpected bytes: got %v, want %v", b, expected)
	}
}

func TestReadErrorEncode(t *testing.T) {
	w := New(func() {}, nil)
	_, err := io.ReadAll(w)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestWriteTo(t *testing.T) {
	s := "test"
	w := New(s, nil)
	buf := new(bytes.Buffer)
	n, err := w.WriteTo(buf)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if n != 7 {
		t.Fatalf("unexpected bytes written: got %d, want %d", n, 7)
	}
}

func TestWriteToErrorEncode(t *testing.T) {
	w := New(func() {}, nil)
	buf := new(bytes.Buffer)
	_, err := w.WriteTo(buf)
	if err == nil {
		t.Fatal("no error")
	}
}
