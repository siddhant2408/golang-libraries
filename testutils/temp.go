package testutils

import (
	"os"
	"testing"

	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/tmpfs"
)

// TempDir is a helper for os.MkdirTemp.
//
// The returned close function deletes the directory.
func TempDir(tb testing.TB, dir string, prefix string) (name string, cl closeutils.F) {
	name, cl, err := tmpfs.Dir(dir, prefix)
	if err != nil {
		FatalErr(tb, err)
	}
	return name, cl
}

// TempFile is a helper for os.CreateTemp.
//
// The returned close function closes and deletes the file.
func TempFile(tb testing.TB, dir string, pattern string) (f *os.File, cl closeutils.F) {
	f, cl, err := tmpfs.File(dir, pattern)
	if err != nil {
		FatalErr(tb, err)
	}
	return f, cl
}
