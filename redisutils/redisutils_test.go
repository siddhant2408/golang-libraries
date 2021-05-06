package redisutils

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/siddhant2408/golang-libraries/redistest"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	redistest.CheckAvailable(t)
	appName := "test"
	c := NewClient(&redis.Options{
		Addr: redistest.GetAddress(),
	}, appName)
	s, err := c.ClientGetName(ctx).Result()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if s != appName {
		t.Fatalf("unexpected name: got %q, want %s", s, appName)
	}
}
