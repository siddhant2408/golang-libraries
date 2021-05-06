package amqptest

import (
	"testing"
)

const testVhost = "test"

func TestCheckAvailable(t *testing.T) {
	CheckAvailable(t)
}

func TestNewConnectionManager(t *testing.T) {
	cm := NewConnectionManager(t, testVhost)
	if cm == nil {
		t.Fatal("nil")
	}
}

func TestVhost(t *testing.T) {
	Vhost(t, testVhost)
}
