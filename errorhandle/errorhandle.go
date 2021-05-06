// Package errorhandle provides a helper function to handle errors.
//
// The helper function:
//  - sends the error to Sentry
//  - add tags to the tracing span
//  - log it
package errorhandle

import (
	"context"
	"net/http"

	raven "github.com/getsentry/raven-go"
	"github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/errorlog"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/goroutine"
	"github.com/siddhant2408/golang-libraries/ravenerrors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

// Handle handles the error.
func Handle(ctx context.Context, myerr error, opts ...Option) {
	if errors.IsIgnored(myerr) {
		return
	}
	cfg := getConfig(opts...)
	if cfg.fatal {
		myerr = ravenerrors.WithSeverity(myerr, raven.FATAL)
	}
	rf := getRavenFunc(cfg)
	ris := getRavenInterfaces(myerr)
	sentryID := rf(myerr, nil, ris...)
	myerr = errors.WithValue(myerr, "sentry.id", sentryID)
	setSentryID(cfg, sentryID)
	setTraceSpanTags(ctx, myerr, sentryID)
	setHTTPHeader(cfg, sentryID)
	lf := getLogFunc(cfg)
	lf(myerr)
}

// HandleDefault calls Handle without options.
func HandleDefault(ctx context.Context, myerr error) {
	Handle(ctx, myerr)
}

// NewHandler returns a new handler function for the given options.
func NewHandler(opts ...Option) func(context.Context, error) {
	return func(ctx context.Context, myerr error) {
		Handle(ctx, myerr, opts...)
	}
}

// Option represents an option for Handle().
type Option func(*config)

type config struct {
	wait       bool
	fatal      bool
	sentryID   *string
	httpHeader httpHeader
}

func getConfig(opts ...Option) *config {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// Wait is an option that blocks until the error is processed.
func Wait() Option {
	return func(cfg *config) {
		cfg.wait = true
	}
}

// Fatal is an option that closes the application after the error is processed.
//
// It also implies Wait().
func Fatal() Option {
	return func(cfg *config) {
		cfg.fatal = true
		Wait()(cfg)
	}
}

// SentryID is an option that allows to get the Sentry event ID.
func SentryID(sp *string) Option {
	return func(cfg *config) {
		cfg.sentryID = sp
	}
}

func setSentryID(cfg *config, sentryID string) {
	if cfg.sentryID != nil {
		*cfg.sentryID = sentryID
	}
}

const (
	httpHeaderSentry = "X-Sentry-Id"
)

type httpHeader interface {
	Set(key string, value string)
}

// HTTPHeader is an option that sets an HTTP header with the Sentry event ID.
func HTTPHeader(h httpHeader) Option {
	return func(cfg *config) {
		cfg.httpHeader = h
	}
}

func setHTTPHeader(cfg *config, sentryID string) {
	if cfg.httpHeader != nil {
		cfg.httpHeader.Set(httpHeaderSentry, sentryID)
	}
}

func getRavenInterfaces(myerr error) []raven.Interface {
	var itfs []raven.Interface
	if hi := getRavenInterfaceHTTPRequest(myerr); hi != nil {
		itfs = append(itfs, hi)
	}
	return itfs
}

func getRavenInterfaceHTTPRequest(myerr error) *raven.Http {
	req := getHTTPServerRequest(myerr)
	if req != nil {
		return raven.NewHttp(req)
	}
	return nil
}

// getHTTPServerRequest is a copy of httperrors.GetServerRequest. (it doesn't require to import this package).
func getHTTPServerRequest(myerr error) *http.Request {
	var werr interface {
		HTTPServerRequest() *http.Request
	}
	ok := errors.As(myerr, &werr)
	if ok {
		return werr.HTTPServerRequest()
	}
	return nil
}

type ravenFunc func(myerr error, tags map[string]string, interfaces ...raven.Interface) (ravenID string)

func getRavenFunc(cfg *config) ravenFunc {
	if cfg.wait {
		return ravenCaptureAndWait
	}
	return ravenCapture
}

func ravenCapture(myerr error, tags map[string]string, interfaces ...raven.Interface) (ravenID string) {
	ravenID, ch := ravenerrors.Capture(myerr, tags, interfaces...)
	consumeRavenError(ch, myerr)
	return ravenID
}

func ravenCaptureAndWait(myerr error, tags map[string]string, interfaces ...raven.Interface) (ravenID string) {
	ravenID, err := ravenerrors.CaptureAndWait(myerr, tags, interfaces...)
	if err != nil {
		handleRavenError(err, myerr)
	}
	return ravenID
}

const ravenErrorsChannelSize = 1000

var ravenErrors = make(chan *ravenError, ravenErrorsChannelSize)

func init() {
	_ = goroutine.Go(consumeRavenErrors)
}

type ravenError struct {
	ch    <-chan error
	myerr error
}

func consumeRavenError(ch <-chan error, myerr error) {
	rerr := &ravenError{
		ch:    ch,
		myerr: myerr,
	}
	select {
	case ravenErrors <- rerr:
	default:
		err := errors.New("Raven errors channel is full")
		handleRavenError(err, myerr)
	}
}

func consumeRavenErrors() {
	for rerr := range ravenErrors {
		err := <-rerr.ch
		if err != nil {
			err = errors.WithStack(err)
			handleRavenError(err, rerr.myerr)
		}
	}
}

func handleRavenError(err error, myerr error) {
	if !filterRavenError(err) {
		return
	}
	err = errors.WithValue(err, "original_error", myerr.Error())
	err = errors.Wrap(err, "raven")
	err = errors.Wrap(err, "errorhandle")
	errorlog.Print(err)
}

const (
	sentryErrorDropWebCrawlers = "raven: got http status 403 - x-sentry-error: Event dropped due to filter: web-crawlers"
)

func filterRavenError(err error) bool {
	// This code compares the error message.
	// It's bad, but there is not other way to do it.
	msg := errors.UnwrapAll(err).Error()
	return msg != sentryErrorDropWebCrawlers
}

const (
	traceSpanTagSentry = "sentry.id"
)

func setTraceSpanTags(ctx context.Context, myerr error, sentryID string) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}
	tracingutils.SetSpanError(span, myerr)
	span.SetTag(traceSpanTagSentry, sentryID)
}

type logFunc func(error)

func getLogFunc(cfg *config) logFunc {
	if cfg.fatal {
		return errorlog.Fatal
	}
	return errorlog.Print
}
