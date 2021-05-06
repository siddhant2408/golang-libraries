package sibhttpapp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestSetGet(t *testing.T) {
	appName := "test"
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	Set(req, appName)
	res, err := Get(req)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if res != appName {
		t.Fatalf("unexpected application name: got %q, want %q", res, appName)
	}
}

func TestGetErrorNotSet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	_, err := Get(req)
	if err == nil {
		t.Fatal("no error")
	}
}
