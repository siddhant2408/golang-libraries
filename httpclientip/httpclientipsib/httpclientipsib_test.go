package httpclientipsib

import (
	"bufio"
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/httpclientrequest"
	"github.com/siddhant2408/golang-libraries/httptestlocal"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestTrusted(t *testing.T) {
	testTrusted(t, "10.0.0.0/8", "Private-Use")
	testTrusted(t, "127.0.0.0/8", "Loopback")
	testTrusted(t, "172.16.0.0/12", "Private-Use")
	testTrusted(t, "192.168.0.0/16", "Private-Use")
	testTrusted(t, "::1/128", "Loopback")
	testTrusted(t, "fc00::/7", "Unique Local")
	testTrustedURLFile(t, "https://www.cloudflare.com/ips-v4")
	testTrustedURLFile(t, "https://www.cloudflare.com/ips-v6")
}

func testTrustedURLFile(t *testing.T, u string) {
	ctx := context.Background()
	ctx = httptestlocal.Allow(ctx)
	res, err := httpclientrequest.Get(ctx, u)
	if err != nil {
		err = errors.Wrap(err, "download file")
		testutils.SkipErr(t, err)
	}
	sc := bufio.NewScanner(bytes.NewReader(res.Body))
	for sc.Scan() {
		testTrusted(t, sc.Text(), u)
	}
	err = sc.Err()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func testTrusted(t *testing.T, v string, name string) {
	for _, t := range trusted {
		if t == v {
			return
		}
	}
	t.Fatalf("not trusted: %s (%s)", v, name)
}

func TestGetter(t *testing.T) {
	for _, tc := range []struct {
		name       string
		addr       string
		xffs       []string
		expectedIP net.IP
	}{
		{
			name:       "Direct",
			addr:       "1.1.1.1:12345",
			expectedIP: net.ParseIP("1.1.1.1"),
		},
		{
			name:       "HAProxy",
			addr:       "172.20.0.10:12345",
			xffs:       []string{"1.1.1.1"},
			expectedIP: net.ParseIP("1.1.1.1"),
		},
		{
			name:       "Cloudflare",
			addr:       "172.20.0.10:12345",
			xffs:       []string{"1.1.1.1, 104.16.1.1"},
			expectedIP: net.ParseIP("1.1.1.1"),
		},
		{
			name:       "AWS",
			addr:       "10.1.2.3:12345",
			xffs:       []string{"1.1.1.1"},
			expectedIP: net.ParseIP("1.1.1.1"),
		},
		{
			name:       "Proxy",
			addr:       "172.20.0.10:12345",
			xffs:       []string{"1.1.1.1, 2.2.2.2, 104.16.1.1"},
			expectedIP: net.ParseIP("2.2.2.2"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
			req.RemoteAddr = tc.addr
			for _, xff := range tc.xffs {
				req.Header.Add("X-Forwarded-For", xff)
			}
			ip, err := Getter.GetClientIP(req)
			if err != nil {
				testutils.FatalErr(t, err)
			}
			if !ip.Equal(tc.expectedIP) {
				t.Fatalf("unexpected result: got %s, want %s", ip, tc.expectedIP)
			}
		})
	}
}
