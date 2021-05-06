// Package jsonreader provides a helper that JSON encodes a value to an io.Reader.
package jsonreader

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/siddhant2408/golang-libraries/errors"
)

type reader struct {
	v   interface{}
	ef  func(*json.Encoder)
	buf io.Reader
}

// New returns a helper that JSON encodes a value to an io.Reader.
//
// The optional ef parameter allows to customize the encoder.
// If nil, NewEncoder() is used.
func New(v interface{}, ef func(*json.Encoder)) interface {
	io.Reader
	io.WriterTo
} {
	return &reader{
		v:  v,
		ef: ef,
	}
}

func (r *reader) Read(p []byte) (int, error) {
	if r.buf == nil {
		buf := new(bytes.Buffer)
		err := r.encode(buf)
		if err != nil {
			return 0, errors.Wrap(err, "encode")
		}
		r.buf = buf
	}
	return r.buf.Read(p)
}

func (r *reader) WriteTo(w io.Writer) (int64, error) {
	cw := &countWriter{
		Writer: w,
	}
	err := r.encode(cw)
	err = errors.Wrap(err, "encode")
	return cw.n, err
}

func (r *reader) encode(w io.Writer) error {
	enc := json.NewEncoder(w)
	if r.ef != nil {
		r.ef(enc)
	}
	return enc.Encode(r.v)
}

type countWriter struct {
	io.Writer
	n int64
}

func (cw *countWriter) Write(p []byte) (int, error) {
	n, err := cw.Writer.Write(p)
	if n > 0 {
		cw.n += int64(n)
	}
	return n, err
}
