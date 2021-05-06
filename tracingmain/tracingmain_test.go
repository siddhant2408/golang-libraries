package tracingmain

import (
	"os"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestGetAgentDebug(t *testing.T) {
	_ = os.Setenv(envDDAgentDebug, "true")
	defer os.Unsetenv(envDDAgentDebug) //nolint:errcheck
	debug, err := getAgentDebug()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if !debug {
		t.Fatalf("unexpected debug: got %t, want %t", debug, true)
	}
}

func TestGetAgentDebugDefault(t *testing.T) {
	debug, err := getAgentDebug()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if debug {
		t.Fatalf("unexpected debug: got %t, want %t", debug, false)
	}
}

func TestGetAgentDebugError(t *testing.T) {
	_ = os.Setenv(envDDAgentDebug, "invalid")
	defer os.Unsetenv(envDDAgentDebug) //nolint:errcheck
	_, err := getAgentDebug()
	if err == nil {
		t.Fatal("no error")
	}
}
