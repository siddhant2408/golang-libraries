package httpclientip

import (
	"context"
	"net"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestContext(t *testing.T) {
	ip1 := net.ParseIP("123.123.123.123")
	ctx := context.Background()
	ctx = SetToContext(ctx, ip1)
	ip2, err := GetFromContext(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if !ip2.Equal(ip1) {
		t.Fatalf("unexpected IP: got %v, want %v", ip2, ip1)
	}
}

func TestContextErrorNotDefined(t *testing.T) {
	ctx := context.Background()
	_, err := GetFromContext(ctx)
	if err == nil {
		t.Fatal("no error")
	}
}
