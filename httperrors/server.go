package httperrors

import (
	"fmt"
	"net/http"

	"github.com/siddhant2408/golang-libraries/errors"
)

// WithServerRequest annotates an error with a server Request.
func WithServerRequest(err error, req *http.Request) error {
	if err == nil {
		return nil
	}
	return &serverRequest{
		err: err,
		req: req,
	}
}

type serverRequest struct {
	err error
	req *http.Request
}

func (err *serverRequest) HTTPServerRequest() *http.Request {
	return err.req
}

func (err *serverRequest) WriteErrorMessage(w errors.Writer, verbose bool) bool {
	writeRequestMessage(w, verbose, err.req, requestTypeServer)
	return true
}

func (err *serverRequest) Error() string                 { return errors.Error(err) }
func (err *serverRequest) Format(s fmt.State, verb rune) { errors.Format(err, s, verb) }
func (err *serverRequest) Unwrap() error                 { return err.err }

// GetServerRequest returns the server Request associated to an error.
//
// If the error is not wrapped with WithServerRequest(), it returns nil.
func GetServerRequest(err error) *http.Request {
	var werr *serverRequest
	ok := errors.As(err, &werr)
	if ok {
		return werr.HTTPServerRequest()
	}
	return nil
}

// WithServerCode annotates an error with an server status code.
//
// See GetServerCodeText() documentation.
func WithServerCode(err error, code int) error {
	if err == nil {
		return nil
	}
	return &serverCode{
		err:  err,
		code: code,
	}
}

type serverCode struct {
	err  error
	code int
}

func (err *serverCode) HTTPServerCode() int {
	return err.code
}

func (err *serverCode) WriteErrorMessage(w errors.Writer, verbose bool) bool {
	fmt.Fprintf(w, "HTTP server code %d", err.code)
	return true
}

func (err *serverCode) Error() string                 { return errors.Error(err) }
func (err *serverCode) Format(s fmt.State, verb rune) { errors.Format(err, s, verb) }
func (err *serverCode) Unwrap() error                 { return err.err }

// GetServerCodeText returns the server status code and text associated to an error.
//
// If the error is wrapped with WithServerCode(), the status code code is the provided one, and the text is the message of the error that was wrapped with WithServerCode().
//
// If the error is not wrapped with WithServerCode(), it returns a 500 status code and a generic message.
// This should be considered the default behavior of all errors, so it's not necessary to call WithServerCode() with http.StatusInternalError/500.
func GetServerCodeText(err error) (code int, text string) {
	var werr *serverCode
	ok := errors.As(err, &werr)
	if ok {
		return werr.HTTPServerCode(), werr.err.Error()
	}
	return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
}
