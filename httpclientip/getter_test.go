package httpclientip

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestNewGetter(t *testing.T) {
	g, err := NewGetter([]string{"1.2.3.4", "1.2.3.4/32"})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if g == nil {
		t.Fatal("nil")
	}
}

func TestNewGetterErrorIP(t *testing.T) {
	g, err := NewGetter([]string{"invalid"})
	if err == nil {
		t.Fatal("no error")
	}
	if g != nil {
		t.Fatal("not nil")
	}
}

func TestNewGetterErrorCIDR(t *testing.T) {
	g, err := NewGetter([]string{"invalid/"})
	if err == nil {
		t.Fatal("no error")
	}
	if g != nil {
		t.Fatal("not nil")
	}
}

func TestGetter(t *testing.T) {
	for _, tc := range []struct {
		name          string
		trusted       []string
		addr          string
		xffs          []string
		expectedIP    net.IP
		expectedError bool
	}{
		{
			name:       "Normal",
			addr:       "1.2.3.4:12345",
			expectedIP: net.ParseIP("1.2.3.4"),
		},
		{
			name:       "NoTrusted",
			addr:       "1.1.1.1:12345",
			xffs:       []string{"2.2.2.2, 3.3.3.3"},
			expectedIP: net.ParseIP("1.1.1.1"),
		},
		{
			name:       "Trusted1",
			trusted:    []string{"1.1.1.1"},
			addr:       "1.1.1.1:12345",
			xffs:       []string{"2.2.2.2"},
			expectedIP: net.ParseIP("2.2.2.2"),
		},
		{
			name:       "Trusted2",
			trusted:    []string{"1.1.1.1", "2.2.2.2"},
			addr:       "1.1.1.1:12345",
			xffs:       []string{"3.3.3.3, 2.2.2.2"},
			expectedIP: net.ParseIP("3.3.3.3"),
		},
		{
			name:       "MultiXFF",
			trusted:    []string{"1.1.1.1", "2.2.2.2"},
			addr:       "1.1.1.1:12345",
			xffs:       []string{"3.3.3.3", "2.2.2.2"},
			expectedIP: net.ParseIP("3.3.3.3"),
		},
		{
			name:       "Proxy",
			trusted:    []string{"1.1.1.1"},
			addr:       "1.1.1.1:12345",
			xffs:       []string{"3.3.3.3, 2.2.2.2"},
			expectedIP: net.ParseIP("2.2.2.2"),
		},
		{
			name:       "Range",
			trusted:    []string{"1.1.1.0/24"},
			addr:       "1.1.1.123:12345",
			xffs:       []string{"2.2.2.2"},
			expectedIP: net.ParseIP("2.2.2.2"),
		},
		{
			name:       "Last",
			trusted:    []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			addr:       "1.1.1.1:12345",
			xffs:       []string{"2.2.2.2, 3.3.3.3"},
			expectedIP: net.ParseIP("2.2.2.2"),
		},
		{
			name:          "ErrorSplitHostPort",
			addr:          "invalid",
			expectedError: true,
		},
		{
			name:          "ErrorParseIP",
			addr:          "invalid:12345",
			expectedError: true,
		},
		{
			name:          "ErrorNoClientIPEmpty",
			expectedError: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewGetter(tc.trusted)
			if err != nil {
				testutils.FatalErr(t, err)
			}
			req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
			req.RemoteAddr = tc.addr
			for _, xff := range tc.xffs {
				req.Header.Add("X-Forwarded-For", xff)
			}
			ip, err := g.GetClientIP(req)
			if err != nil {
				if tc.expectedError {
					return
				}
				testutils.FatalErr(t, err)
			}
			if tc.expectedError {
				t.Fatal("no error")
			}
			if !ip.Equal(tc.expectedIP) {
				t.Fatalf("unexpected result: got %s, want %s", ip, tc.expectedIP)
			}
		})
	}
}
