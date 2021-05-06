// Package spewutils configures the spew package for our internal usage.
package spewutils

import (
	"bytes"
	"io"

	"github.com/davecgh/go-spew/spew"
	"github.com/siddhant2408/golang-libraries/bufpool"
	"github.com/siddhant2408/golang-libraries/errors"
)

func init() {
	spew.Config.Indent = "\t"
	spew.Config.DisablePointerAddresses = true
	spew.Config.DisableCapacities = true
	spew.Config.SortKeys = true
	spew.Config.SpewKeys = false
	errors.ValueWriter = writeValueWithoutNewline
}

func writeValueWithoutNewline(w io.Writer, v interface{}) {
	buf := bufpool.Get()
	defer bufpool.Put(buf)
	spew.Fdump(buf, v)
	b := buf.Bytes()
	b = bytes.TrimSuffix(b, []byte("\n"))
	_, _ = w.Write(b)
}
