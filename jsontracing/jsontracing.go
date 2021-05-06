// Package jsontracing provides JSON tracing related utilities.
package jsontracing

import (
	"context"
	"encoding/json"
	"io"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

// Marshal is a helper for json.Marshal.
func Marshal(ctx context.Context, v interface{}) (b []byte, err error) {
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "json.marshal", &err)
	defer spanFinish()
	b, err = json.Marshal(v)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	span.SetTag("json.size", len(b))
	return b, nil
}

// Unmarshal is a helper for json.Unmarshal.
func Unmarshal(ctx context.Context, data []byte, v interface{}) (err error) {
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "json.unmarshal", &err)
	defer spanFinish()
	span.SetTag("json.size", len(data))
	err = json.Unmarshal(data, v)
	err = errors.Wrap(err, "")
	return err
}

// Encoder is a wrapper for json.Encoder.
type Encoder struct {
	*json.Encoder
}

// NewEncoder returns a new Encoder.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		Encoder: json.NewEncoder(w),
	}
}

// Encode is a helper for json.Encoder.Encode.
func (e *Encoder) Encode(ctx context.Context, v interface{}) (err error) {
	_, spanFinish := tracingutils.StartChildSpan(&ctx, "json.encoder", &err)
	defer spanFinish()
	err = e.Encoder.Encode(v)
	err = errors.Wrap(err, "")
	return err
}

// Decoder is a wrapper for json.Decoder.
type Decoder struct {
	*json.Decoder
}

// NewDecoder returns a new Decoder.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		Decoder: json.NewDecoder(r),
	}
}

// Decode is a helper for json.Decoder.Decode.
func (d *Decoder) Decode(ctx context.Context, v interface{}) (err error) {
	_, spanFinish := tracingutils.StartChildSpan(&ctx, "json.decoder", &err)
	defer spanFinish()
	err = d.Decoder.Decode(v)
	err = errors.Wrap(err, "")
	return err
}
