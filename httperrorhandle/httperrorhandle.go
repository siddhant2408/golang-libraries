// Package httperrorhandle provides error handling for HTTP handler.
package httperrorhandle

import (
	"context"
	"net/http"

	"github.com/siddhant2408/golang-libraries/errorhandle"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/httperrors"
)

// Handle handles an error for an HTTP handler.
//
// It returns the code/text for the given error.
// If the code is in the 5XX range, it calls errorhandle and sets the Sentry header in the response.
//
// It doesn't write a response.
func Handle(ctx context.Context, w http.ResponseWriter, req *http.Request, err error) *HandleResult {
	res := new(HandleResult)
	res.Code, res.Text = httperrors.GetServerCodeText(err)
	if res.Code >= 500 && res.Code < 600 {
		err = httperrors.WithServerRequest(err, req)
		errorhandle.Handle(
			ctx,
			err,
			errorhandle.HTTPHeader(w.Header()),
			errorhandle.SentryID(&res.SentryID),
		)
	}
	return res
}

// HandleResult is the result of Handle.
type HandleResult struct {
	Code     int
	Text     string
	SentryID string
}

// Serve serves an HTTP request and handles the returned error.
func Serve(w http.ResponseWriter, req *http.Request, f Func, onErr OnErrFunc) {
	ctx := req.Context()
	err := f(ctx, w, req)
	if err == nil {
		return
	}
	err = errors.Wrap(err, "HTTP serve")
	hres := Handle(ctx, w, req, err)
	onErr(ctx, w, req, err, hres)
}

// Func is called by Server in order to server the request.
// If an error is returned, it must not write any response.
type Func func(ctx context.Context, w http.ResponseWriter, req *http.Request) error

// OnErrFunc Is called by Serve if Func returns an error.
// It must return an error response.
type OnErrFunc func(ctx context.Context, w http.ResponseWriter, req *http.Request, err error, hres *HandleResult)
