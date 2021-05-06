// Package httputils provides HTTP related utilities.
package httputils

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/iotracing"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

// ReadRequestBody reads an HTTP request body.
func ReadRequestBody(ctx context.Context, req *http.Request) (b []byte, err error) {
	_, spanFinish := tracingutils.StartChildSpan(&ctx, "http.read_request_body", &err)
	defer spanFinish()
	b, err = iotracing.ReadAll(ctx, req.Body)
	err = errors.Wrap(err, "read all")
	return b, err
}

// CopyRequestBody copies an HTTP request body.
func CopyRequestBody(ctx context.Context, req *http.Request, w io.Writer) (written int64, err error) {
	_, spanFinish := tracingutils.StartChildSpan(&ctx, "http.copy_request_body", &err)
	defer spanFinish()
	written, err = iotracing.Copy(ctx, w, req.Body)
	err = errors.Wrap(err, "copy")
	return written, err
}

// WriteResponse writes an HTTP response.
//
// Write errors are ignored.
func WriteResponse(ctx context.Context, w http.ResponseWriter, code int, data []byte) {
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "http.write_response", nil)
	defer spanFinish()
	l := len(data)
	span.SetTag("http.status_code", code)
	span.SetTag("http.content_length", l)
	setHeaderContentLength(w.Header(), l)
	w.WriteHeader(code)
	if l > 0 {
		_, _ = iotracing.Write(ctx, w, data)
	}
}

// WriteResponseText writes a text HTTP response.
//
// Write errors are ignored.
func WriteResponseText(ctx context.Context, w http.ResponseWriter, code int, data string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	WriteResponse(ctx, w, code, []byte(data))
}

// CopyResponse copies an HTTP response.
//
// Copy errors are ignored.
func CopyResponse(ctx context.Context, w http.ResponseWriter, code int, r io.Reader) {
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "http.copy_response", nil)
	defer spanFinish()
	span.SetTag("http.status_code", code)
	if r, ok := r.(interface {
		Len() int
	}); ok {
		l := r.Len()
		span.SetTag("http.content_length", l)
		setHeaderContentLength(w.Header(), l)
	}
	w.WriteHeader(code)
	_, _ = iotracing.Copy(ctx, w, r)
}

func setHeaderContentLength(hd http.Header, l int) {
	hd.Set("Content-Length", strconv.Itoa(l))
}
