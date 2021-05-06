package kafkatest

import (
	"testing"
)

func TestCheckAvailable(t *testing.T) {
	CheckAvailable(t)
}

func TestTopic(t *testing.T) {
	_ = Topic(t)
}
