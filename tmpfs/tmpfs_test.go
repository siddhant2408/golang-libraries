package tmpfs_test

import (
	"os"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
	. "github.com/siddhant2408/golang-libraries/tmpfs"
)

func TestDir(t *testing.T) {
	name, cl, err := Dir("", "")
	if err != nil {
		testutils.FatalErr(t, err)
	}
	defer cl()
	f, err := os.Open(name)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	fi, err := f.Stat()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	d := fi.IsDir()
	if !d {
		t.Fatal("not dir")
	}
}

func TestFile(t *testing.T) {
	f, cl, err := File("", "")
	if err != nil {
		testutils.FatalErr(t, err)
	}
	defer cl()
	_, err = f.Stat()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
