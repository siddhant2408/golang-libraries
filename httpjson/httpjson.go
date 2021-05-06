// Package httpjson provides HTTP JSON related utilities.
package httpjson

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/siddhant2408/golang-libraries/bufpool"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/httputils"
	"github.com/siddhant2408/golang-libraries/jsontracing"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

// DecodeRequestBody reads the request body and JSON decode it.
func DecodeRequestBody(ctx context.Context, req *http.Request, v interface{}, oo ...RequestOption) (err error) {
	_, spanFinish := tracingutils.StartChildSpan(&ctx, "httpjson.read_request", &err)
	defer spanFinish()
	opts := getRequestOptions(oo...)
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	_, err = httputils.CopyRequestBody(ctx, req, buf)
	if err != nil {
		return errors.Wrap(err, "copy request body")
	}
	dec := jsontracing.NewDecoder(buf)
	if opts.decoder != nil {
		opts.decoder(dec.Decoder)
	}
	err = dec.Decode(ctx, v)
	if err != nil {
		return errors.Wrap(err, "decode")
	}
	err = checkDecoderRemainingData(dec)
	if err != nil {
		return errors.Wrap(err, "remaining data")
	}
	return nil
}

func checkDecoderRemainingData(dec *jsontracing.Decoder) error {
	// Check if there is more data after the first JSON object.
	// We want to return an error if there is a second JSON object after the first one, or if there is invalid JSON data.
	// However we don't want to return an error if there are some trailing spaces or new lines characters.
	tok, err := dec.Token()
	if err == nil {
		// More valid data.
		return errors.Newf("token: %q", tok)
	}
	if errors.Is(err, io.EOF) {
		// No more data.
		return nil
	}
	// More invalid data.
	return errors.Wrap(err, "invalid")
}

func getRequestOptions(oo ...RequestOption) *requestOptions {
	opts := &requestOptions{}
	for _, o := range oo {
		o(opts)
	}
	return opts
}

type requestOptions struct {
	decoder func(*json.Decoder)
}

// RequestOption represents request option.
type RequestOption func(*requestOptions)

// RequestDecoder configures a request decoder.
func RequestDecoder(f func(*json.Decoder)) RequestOption {
	return func(o *requestOptions) {
		o.decoder = f
	}
}

// WriteResponse writes a JSON encoded response.
// If the JSON marshaling returns an error, it's returned without writing any response.
//
// Write errors are ignored.
func WriteResponse(ctx context.Context, w http.ResponseWriter, code int, data interface{}, oo ...ResponseOption) (err error) {
	_, spanFinish := tracingutils.StartChildSpan(&ctx, "httpjson.write_response", &err)
	defer spanFinish()
	opts := getResponseOptions(oo...)
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	enc := jsontracing.NewEncoder(buf)
	if opts.encoder != nil {
		opts.encoder(enc.Encoder)
	}
	err = enc.Encode(ctx, data)
	if err != nil {
		return errors.Wrap(err, "encode")
	}
	hd := w.Header()
	hd.Set("Content-Type", "application/json")
	if opts.header != nil {
		opts.header(hd)
	}
	httputils.CopyResponse(ctx, w, code, buf)
	return nil
}

func getResponseOptions(oo ...ResponseOption) *responseOptions {
	opts := &responseOptions{}
	for _, o := range oo {
		o(opts)
	}
	return opts
}

type responseOptions struct {
	encoder func(*json.Encoder)
	header  func(http.Header)
}

// ResponseOption represents response option.
type ResponseOption func(*responseOptions)

// ResponseEncoder configures a response encoder.
func ResponseEncoder(f func(*json.Encoder)) ResponseOption {
	return func(opts *responseOptions) {
		opts.encoder = f
	}
}

// ResponseHeader configures a response header.
func ResponseHeader(f func(http.Header)) ResponseOption {
	return func(opts *responseOptions) {
		opts.header = f
	}
}

var bufPool = bufpool.New()

func init() {
	bufPool.MaxCap = 1 << 20 // 1 MiB
}
