// Package iotracing provides IO tracing related utilities.
package iotracing

import (
	"context"
	"io"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

// Read is a helper for io.Reader.Read.
func Read(ctx context.Context, r io.Reader, p []byte) (n int, err error) {
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "io.read", &err)
	defer spanFinish()
	n, err = r.Read(p)
	err = errors.Wrap(err, "")
	span.SetTag("bytes", n)
	return n, err
}

// Write is a helper for io.Writer.Write.
func Write(ctx context.Context, w io.Writer, p []byte) (n int, err error) {
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "io.write", &err)
	defer spanFinish()
	span.SetTag("bytes", len(p))
	n, err = w.Write(p)
	err = errors.Wrap(err, "")
	return n, err
}

// Copy is a helper for io.Copy.
func Copy(ctx context.Context, dst io.Writer, src io.Reader) (written int64, err error) {
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "io.copy", &err)
	defer spanFinish()
	written, err = io.Copy(dst, src)
	err = errors.Wrap(err, "")
	span.SetTag("bytes", written)
	return written, err
}

// ReadAll is a helper for io.ReadAll.
func ReadAll(ctx context.Context, r io.Reader) (b []byte, err error) {
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "io.read_all", &err)
	defer spanFinish()
	b, err = io.ReadAll(r)
	err = errors.Wrap(err, "")
	span.SetTag("bytes", len(b))
	return b, err
}
