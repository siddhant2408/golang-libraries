package debugutils

import (
	"encoding/json"
	"expvar"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestHTTPHandler(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost/debug", nil)
	httpHandler(w, req)
	w.Flush()
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.Len() == 0 {
		t.Fatal("empty body")
	}
}

func TestExpvarEnvVars(t *testing.T) {
	v := expvar.Get("envvars")
	var m map[string]string
	err := json.Unmarshal([]byte(v.String()), &m)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if len(m) == 0 {
		t.Fatal("no variables")
	}
}

func TestExpvarGoroutines(t *testing.T) {
	v := expvar.Get("goroutines")
	_, err := strconv.Atoi(v.String())
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestExpvarPID(t *testing.T) {
	v := expvar.Get("pid")
	_, err := strconv.Atoi(v.String())
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestExpvarUptime(t *testing.T) {
	v := expvar.Get("uptime")
	_, err := strconv.ParseFloat(v.String(), 64)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
