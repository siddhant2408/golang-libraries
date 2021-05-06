// Package httpclientrequest provides a helper to do HTTP client request.
package httpclientrequest

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/httperrors"
	"github.com/siddhant2408/golang-libraries/iotracing"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

// Do executes an HTTP request.
// It downloads the response body to a []byte and closes it.
//
// Default options:
//  - client: http.DefaultClient
//  - timeout: 30 sec
//  - max body size: 10 MiB
//  - copy body: nil
//  - status check: 2XX
//
// The context associated to the request is ignored.
// Use the one passed in parameters instead.
//
// The returned response and body may be defined even if an error is returned.
func Do(ctx context.Context, req *http.Request, opts ...Option) (res Result, err error) {
	_, spanFinish := startTraceSpan(&ctx, "do", &err)
	defer spanFinish()
	o := getOptions(opts...)
	if o.timeout > 0 {
		var cancelTimeout context.CancelFunc
		ctx, cancelTimeout = context.WithTimeout(ctx, o.timeout)
		defer cancelTimeout()
	}
	resp, err := doRequest(ctx, o.client, req)
	res.Response = resp // This avoids the bodyclose lint warning.
	if err != nil {
		err = httperrors.WithClientRequest(err, req)
		err = errors.Wrap(err, "request")
		return res, err
	}
	res.Body, err = processResponse(ctx, resp, o)
	if err != nil {
		err = httperrors.WithClientResponse(err, &httperrors.ClientResponse{
			Response: resp,
			Body:     res.Body,
		})
		err = httperrors.WithClientRequest(err, req)
		err = errors.Wrap(err, "response")
		return res, err
	}
	return res, nil
}

func doRequest(ctx context.Context, clt *http.Client, req *http.Request) (resp *http.Response, err error) {
	_, spanFinish := startTraceSpan(&ctx, "request", &err)
	defer spanFinish()
	req = req.WithContext(ctx)
	resp, err = clt.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client do")
	}
	return resp, nil
}

func processResponse(ctx context.Context, resp *http.Response, o *options) ([]byte, error) {
	applyMaxBodySize(resp, o)
	body, err := readResponseBody(ctx, resp, o)
	if err != nil {
		return body, errors.Wrap(err, "read body")
	}
	err = checkStatus(resp, o.status)
	if err != nil {
		return body, errors.Wrap(err, "status")
	}
	return body, nil
}

func applyMaxBodySize(resp *http.Response, o *options) {
	if o.maxBodySize > 0 {
		resp.Body = &readCloser{
			Reader: newMaxBytesReader(resp.Body, o.maxBodySize),
			Closer: resp.Body,
		}
	}
}

func readResponseBody(ctx context.Context, resp *http.Response, o *options) (body []byte, err error) {
	var bodySize int64
	span, spanFinish := startTraceSpan(&ctx, "read_body", &err)
	defer func() {
		// Ensure that we can't interact with the response body.
		_ = resp.Body.Close()
		resp.Body = &panicResponseBody{}
		span.SetTag("http.response.body.size", bodySize)
		spanFinish()
	}()
	if o.copyBody != nil {
		bodySize, err = iotracing.Copy(ctx, o.copyBody, resp.Body)
		err = errors.Wrap(err, "copy")
	} else {
		body, err = io.ReadAll(resp.Body)
		bodySize = int64(len(body))
		err = errors.Wrap(err, "read all")
	}
	return body, err
}

func checkStatus(resp *http.Response, f func(int) error) error {
	if f != nil {
		return f(resp.StatusCode)
	}
	return nil
}

func startTraceSpan(pctx *context.Context, operationName string, perr *error) (opentracing.Span, closeutils.F) {
	return tracingutils.StartChildSpan(pctx, "httpclientrequest."+operationName, perr)
}

// Option represents an option.
type Option func(*options)

type options struct {
	client      *http.Client
	timeout     time.Duration
	maxBodySize int64
	copyBody    io.Writer
	status      func(int) error
}

const (
	timeoutDefault     = 30 * time.Second
	maxBodySizeDefault = 10 * 1 << 20 // 10 MiB
)

func newOptions() *options {
	return &options{
		client:      http.DefaultClient,
		timeout:     timeoutDefault,
		maxBodySize: maxBodySizeDefault,
		status:      status2XX,
	}
}

func getOptions(opts ...Option) *options {
	o := newOptions()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Client defines the client.
func Client(c *http.Client) Option {
	return func(o *options) {
		o.client = c
	}
}

// Timeout defines the timeout.
//
// 0 means no timeout.
func Timeout(to time.Duration) Option {
	return func(o *options) {
		o.timeout = to
	}
}

// MaxBodySize defines the maximum response body size.
//
// 0 means no maximum size.
func MaxBodySize(max int64) Option {
	return func(o *options) {
		o.maxBodySize = max
	}
}

// CopyBody copies the body to a io.Writer.
//
// If nil, the body is read to Result.Body.
func CopyBody(w io.Writer) Option {
	return func(o *options) {
		o.copyBody = w
	}
}

// StatusFunc checks the status with a function.
//
// Nil disables the status check.
func StatusFunc(f func(int) error) Option {
	return func(o *options) {
		o.status = f
	}
}

// Status checks the status with a single value.
func Status(st int) Option {
	return StatusFunc(func(s int) error {
		if s != st {
			return errors.Newf("got %d, want %d", s, st)
		}
		return nil
	})
}

// Status2XX checks the status with 2XX.
func Status2XX() Option {
	return StatusFunc(status2XX)
}

func status2XX(s int) error {
	if s < 200 || s >= 300 {
		return errors.Newf("got %d, want 2XX", s)
	}
	return nil
}

type readCloser struct {
	io.Reader
	io.Closer
}

// maxBytesReader is a modified copy of net/http.MaxBytesReader.
// It doesn't require a ResponseWriter and the error message is more generic.
type maxBytesReader struct {
	r   io.Reader
	max int64
	n   int64
	err error
}

func newMaxBytesReader(r io.Reader, max int64) io.Reader {
	return &maxBytesReader{
		r:   r,
		max: max,
		n:   max,
	}
}

func (r *maxBytesReader) Read(p []byte) (n int, err error) {
	if r.err != nil {
		return 0, r.err
	}
	if len(p) == 0 {
		return 0, nil
	}
	if int64(len(p)) > r.n+1 {
		p = p[:r.n+1]
	}
	n, err = r.r.Read(p)
	if int64(n) <= r.n {
		r.n -= int64(n)
		r.err = err
		return n, err
	}
	n = int(r.n)
	r.n = 0
	r.err = errors.Newf("more than %d bytes", r.max)
	return n, r.err
}

// panicResponseBody is an implementation for Response.Body that panics for each call.
// It prevents the developer from using the body of the response.
//
// It is fine to panic here, because the developer is doing a mistake.
type panicResponseBody struct{}

const panicResponseBodyMsg = "the response body is managed by httpclientrequest.Do()"

func (prb *panicResponseBody) Read(p []byte) (n int, err error) {
	panic(panicResponseBodyMsg)
}

func (prb *panicResponseBody) Close() error {
	panic(panicResponseBodyMsg)
}

// Result represents the result of an HTTP request.
type Result struct {
	Response *http.Response // The sub-field Response.Body must not be used.
	Body     []byte         // Contains the response body if the option "copy body" is nil.
}

// Get is a helper for http.Get.
func Get(ctx context.Context, u string, opts ...Option) (res Result, err error) {
	return helper(ctx, http.MethodGet, u, "", nil, opts...)
}

// Head is a helper for http.Head.
func Head(ctx context.Context, u string, opts ...Option) (res Result, err error) {
	return helper(ctx, http.MethodHead, u, "", nil, opts...)
}

// Post is a helper for http.Post.
func Post(ctx context.Context, u string, contentType string, body io.Reader, opts ...Option) (res Result, err error) {
	return helper(ctx, http.MethodPost, u, contentType, body, opts...)
}

// PostForm is a helper for http.PostForm.
func PostForm(ctx context.Context, u string, data url.Values, opts ...Option) (res Result, err error) {
	return Post(ctx, u, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), opts...)
}

func helper(ctx context.Context, met string, u string, contentType string, body io.Reader, opts ...Option) (res Result, err error) {
	req, err := http.NewRequest(met, u, body)
	if err != nil {
		return res, errors.Wrap(err, "new request")
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return Do(ctx, req, opts...)
}
