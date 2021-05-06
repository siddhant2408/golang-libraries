// Package httptemplate provides a helper to write template to HTTP response.
package httptemplate

import (
	"context"
	"io"
	"net/http"

	"github.com/siddhant2408/golang-libraries/bufpool"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/httputils"
	"github.com/siddhant2408/golang-libraries/templatetracing"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

// WriteResponse writes an HTTP template response.
// If the template execution returns an error, it's returned without writing any response.
//
// The optional hdf function allows to customize the headers.
//
// Write errors are ignored.
func WriteResponse(ctx context.Context, w http.ResponseWriter, hdf func(http.Header), code int, tmpl Template, data interface{}) (err error) {
	_, spanFinish := tracingutils.StartChildSpan(&ctx, "httptemplate.write_response", &err)
	defer spanFinish()
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	err = templatetracing.Execute(ctx, buf, tmpl, data)
	if err != nil {
		return errors.Wrap(err, "execute")
	}
	if hdf != nil {
		hdf(w.Header())
	}
	httputils.CopyResponse(ctx, w, code, buf)
	return nil
}

// Template represents an HTML or text template.
type Template interface {
	Execute(w io.Writer, data interface{}) error
	Name() string
}

var bufPool = bufpool.New()

func init() {
	bufPool.MaxCap = 1 << 20 // 1 MiB
}
