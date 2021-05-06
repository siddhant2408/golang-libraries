package iotracing

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestRead(t *testing.T) {
	ctx := context.Background()
	r := bytes.NewReader([]byte("test"))
	p := make([]byte, 10)
	n, err := Read(ctx, r, p)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if n != 4 {
		t.Fatalf("unexpected count: got %d, want %d", n, 4)
	}
}

func TestWrite(t *testing.T) {
	ctx := context.Background()
	w := io.Discard
	p := []byte("test")
	n, err := Write(ctx, w, p)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if n != 4 {
		t.Fatalf("unexpected count: got %d, want %d", n, 4)
	}
}

func TestCopy(t *testing.T) {
	ctx := context.Background()
	src := bytes.NewReader([]byte("test"))
	dst := io.Discard
	written, err := Copy(ctx, dst, src)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if written != 4 {
		t.Fatalf("unexpected count: got %d, want %d", written, 4)
	}
}

func TestReadAll(t *testing.T) {
	ctx := context.Background()
	b1 := []byte("test")
	r := bytes.NewReader(b1)
	b2, err := ReadAll(ctx, r)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if !bytes.Equal(b2, b1) {
		t.Fatalf("unexpected result: got %v, want %v", b2, b1)
	}
}
