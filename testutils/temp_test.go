package testutils

import (
	"os"
	"testing"
)

func TestTempDir(t *testing.T) {
	name, cl := TempDir(t, "", "")
	defer cl()
	f, err := os.Open(name)
	if err != nil {
		FatalErr(t, err)
	}
	fi, err := f.Stat()
	if err != nil {
		FatalErr(t, err)
	}
	d := fi.IsDir()
	if !d {
		t.Fatal("not dir")
	}
}

func TestTempFile(t *testing.T) {
	f, cl := TempFile(t, "", "")
	defer cl()
	_, err := f.Stat()
	if err != nil {
		FatalErr(t, err)
	}
}
