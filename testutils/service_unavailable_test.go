package testutils

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
)

func TestHandleUnavailable(t *testing.T) {
	envVar := fmt.Sprintf("testServiceSkip:%d", rand.Int63())
	err := os.Setenv(envVar, "true")
	if err != nil {
		FatalErr(t, err)
	}
	defer os.Unsetenv(envVar) //nolint:errcheck
	myErr := errors.New("error")
	HandleUnavailable(t, envVar, myErr)
}

func TestGetBoolEnvVarsFound(t *testing.T) {
	envVar := fmt.Sprintf("testServiceSkip:%d", rand.Int63())
	err := os.Setenv(envVar, "true")
	if err != nil {
		FatalErr(t, err)
	}
	defer os.Unsetenv(envVar) //nolint:errcheck
	value, ok, err := getBoolEnvVars(envVar)
	if err != nil {
		FatalErr(t, err)
	}
	if !ok {
		t.Fatal("not found")
	}
	if !value {
		t.Fatal(false)
	}
}

func TestGetBoolEnvVarsFoundSecond(t *testing.T) {
	envVar1 := fmt.Sprintf("testServiceSkip:%d", rand.Int63())
	envVar2 := fmt.Sprintf("testServiceSkip:%d", rand.Int63())
	err := os.Setenv(envVar2, "true")
	if err != nil {
		FatalErr(t, err)
	}
	defer os.Unsetenv(envVar2) //nolint:errcheck
	value, ok, err := getBoolEnvVars(envVar1, envVar2)
	if err != nil {
		FatalErr(t, err)
	}
	if !ok {
		t.Fatal("not found")
	}
	if !value {
		t.Fatal(false)
	}
}

func TestGetBoolEnvVarsNotFound(t *testing.T) {
	envVar := fmt.Sprintf("testServiceSkip:%d", rand.Int63())
	_, ok, err := getBoolEnvVars(envVar)
	if err != nil {
		FatalErr(t, err)
	}
	if ok {
		t.Fatal("found")
	}
}

func TestGetBoolEnvVarsErrorParseBool(t *testing.T) {
	envVar := fmt.Sprintf("testServiceSkip:%d", rand.Int63())
	err := os.Setenv(envVar, "invalid")
	if err != nil {
		FatalErr(t, err)
	}
	defer os.Unsetenv(envVar) //nolint:errcheck
	_, _, err = getBoolEnvVars(envVar)
	if err == nil {
		t.Fatal("no error")
	}
}
