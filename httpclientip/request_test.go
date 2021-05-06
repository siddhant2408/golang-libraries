package httpclientip

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestGetFromRequest(t *testing.T) {
	ip1 := net.ParseIP("123.123.123.123")
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	ctx := req.Context()
	ctx = SetToContext(ctx, ip1)
	req = req.WithContext(ctx)
	ip2, err := GetFromRequest(req)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if !ip2.Equal(ip1) {
		t.Fatalf("unexpected IP: got %v, want %v", ip2, ip1)
	}
}
