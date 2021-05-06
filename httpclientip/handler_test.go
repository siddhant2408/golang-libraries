package httpclientip

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestHandler(t *testing.T) {
	ip := net.ParseIP("1.2.3.4")
	g, err := NewGetter(nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	h := &Handler{
		Handler: http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
			ipCtx, err := GetFromContext(req.Context())
			if err != nil {
				testutils.FatalErr(t, err)
			}
			if !ipCtx.Equal(ip) {
				t.Fatalf("unexpected IP: got %v, want %v", ipCtx, ip)
			}
		}),
		Getter: g,
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.RemoteAddr = ip.String() + ":12345"
	h.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusOK {
		t.Fatalf("unexpoected status: got %d, want %d", w.Code, http.StatusOK)
	}
}
