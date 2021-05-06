// Package templatetracing provides tracing for template.
package templatetracing

import (
	"context"
	"io"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

// Execute executes a templates and traces the execution.
func Execute(ctx context.Context, w io.Writer, tmpl Template, data interface{}) (err error) {
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "template.execute", &err)
	defer spanFinish()
	span.SetTag("template.name", tmpl.Name())
	err = tmpl.Execute(w, data)
	return errors.Wrap(err, "")
}

// Template represents a template.
type Template interface {
	Execute(w io.Writer, data interface{}) error
	Name() string
}
