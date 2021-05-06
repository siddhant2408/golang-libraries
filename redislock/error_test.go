package redislock

import (
	"testing"
)

func TestErrorString(t *testing.T) {
	msg := ErrNotObtained.Error()
	expected := string(ErrNotObtained)
	if msg != expected {
		t.Fatalf("unexpected message: got %q, want %q", msg, expected)
	}
}
