package mainutils

import (
	"testing"

	"github.com/siddhant2408/golang-libraries/envutils"
	"github.com/siddhant2408/golang-libraries/testutils"
)

var testConfig = Config{
	AppName: "test",
	Version: "1.2.3",
	Env:     envutils.Testing,
}

func TestConfigValidate(t *testing.T) {
	c := testConfig
	err := c.validate()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestConfigValidateErrorAppName(t *testing.T) {
	c := testConfig
	c.AppName = ""
	err := c.validate()
	if err == nil {
		t.Fatal("no error")
	}
}

func TestConfigValidateErrorVersion(t *testing.T) {
	c := testConfig
	c.Version = ""
	err := c.validate()
	if err == nil {
		t.Fatal("no error")
	}
}

func TestConfigValidateErrorEnv(t *testing.T) {
	c := testConfig
	c.Env = envutils.Env("invalid")
	err := c.validate()
	if err == nil {
		t.Fatal("no error")
	}
}
