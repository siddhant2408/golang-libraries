package envutils

import (
	"flag"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestSet(t *testing.T) {
	var e Env
	err := e.Set(string(Testing))
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if e != Testing {
		t.Fatalf("unexpected environment: got %q, want %q", e, Testing)
	}
}

func TestSetError(t *testing.T) {
	var e Env
	err := e.Set("invalid")
	if err == nil {
		t.Fatal("no error")
	}
	if e != Env("") {
		t.Fatalf("unexpected environment: got %q, want %q", e, Env(""))
	}
}

func TestCheck(t *testing.T) {
	for _, env := range []Env{Testing, Development, Staging, Production} {
		err := Check(env)
		if err != nil {
			testutils.FatalErr(t, err)
		}
	}
}

func TestCheckError(t *testing.T) {
	err := Check("invalid")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSetFlag(t *testing.T) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Usage = func() {}
	e := Testing
	SetFlag(fs, &e)
	err := fs.Parse([]string{"-environment=production"})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if e != Production {
		t.Fatalf("unexpected environment: got %q, want %q", e, Production)
	}
}

func TestSetFlagDefault(t *testing.T) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Usage = func() {}
	e := Testing
	SetFlag(fs, &e)
	err := fs.Parse(nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if e != Testing {
		t.Fatalf("unexpected environment: got %q, want %q", e, Testing)
	}
}

func TestSetFlagError(t *testing.T) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Usage = func() {}
	e := Testing
	SetFlag(fs, &e)
	err := fs.Parse([]string{"-environment=invalid"})
	if err == nil {
		t.Fatal("no error")
	}
}
