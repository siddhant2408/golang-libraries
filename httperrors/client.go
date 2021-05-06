package httperrors

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/strconvio"
)

// WithClientRequest annotates an error with a client Request.
func WithClientRequest(err error, req *http.Request) error {
	if err == nil {
		return nil
	}
	return &clientRequest{
		err: err,
		req: req,
	}
}

type clientRequest struct {
	err error
	req *http.Request
}

func (err *clientRequest) HTTPClientRequest() *http.Request {
	return err.req
}

func (err *clientRequest) WriteErrorMessage(w errors.Writer, verbose bool) bool {
	writeRequestMessage(w, verbose, err.req, requestTypeClient)
	return true
}

func (err *clientRequest) Error() string                 { return errors.Error(err) }
func (err *clientRequest) Format(s fmt.State, verb rune) { errors.Format(err, s, verb) }
func (err *clientRequest) Unwrap() error                 { return err.err }

// GetClientRequest returns the client Request associated to an error.
//
// If the error is not wrapped with WithClientRequest(), it returns nil.
func GetClientRequest(err error) *http.Request {
	var werr *clientRequest
	ok := errors.As(err, &werr)
	if ok {
		return werr.HTTPClientRequest()
	}
	return nil
}

const (
	clientResponseBodyMaxSizeString = 4096
	clientResponseBodyMaxSizeBytes  = 1024
)

// WithClientResponse annotates an error with a client response.
func WithClientResponse(err error, value *ClientResponse) error {
	if err == nil {
		return nil
	}
	err = &clientResponse{
		err:   err,
		value: value,
	}
	return err
}

type clientResponse struct {
	err   error
	value *ClientResponse
}

func (err *clientResponse) HTTPClientResponse() *ClientResponse {
	return err.value
}

func (err *clientResponse) WriteErrorMessage(w errors.Writer, verbose bool) bool {
	if !verbose {
		_, _ = w.WriteString("HTTP client response ")
		_, _ = w.WriteString(err.value.Response.Status)
		return true
	}
	_, _ = w.WriteString("HTTP client response\n")
	_, _ = w.WriteString("	status: ")
	_, _ = w.WriteString(err.value.Response.Status)
	_, _ = w.WriteString("\n")
	_, _ = w.WriteString("	proto: ")
	_, _ = w.WriteString(err.value.Response.Proto)
	_, _ = w.WriteString("\n")
	if err.value.Response.ContentLength > 0 {
		_, _ = w.WriteString("	content length: ")
		_, _ = strconvio.WriteInt(w, err.value.Response.ContentLength, 10)
		_, _ = w.WriteString("\n")
	}
	writeHeaderMessage(w, err.value.Response.Header)
	err.writeMessageBody(w)
	return true
}

func (err *clientResponse) writeMessageBody(w errors.Writer) {
	if len(err.value.Body) == 0 {
		return
	}
	_, _ = w.WriteString("	body:\n")
	if utf8.Valid(err.value.Body) {
		err.writeMessageBodyString(w)
	} else {
		err.writeMessageBodyBytes(w)
	}
}

func (err *clientResponse) writeMessageBodyString(w errors.Writer) {
	body, truncated := err.truncateBody(clientResponseBodyMaxSizeString)
	_, _ = w.WriteString("================ begin ================\n")
	_, _ = w.Write(body)
	_, _ = w.WriteString("\n")
	err.writeMessageBodyTruncated(w, truncated, clientResponseBodyMaxSizeString)
	_, _ = w.WriteString("================= end =================\n")
}

func (err *clientResponse) writeMessageBodyBytes(w errors.Writer) {
	body, truncated := err.truncateBody(clientResponseBodyMaxSizeBytes)
	d := hex.Dumper(w)
	_, _ = d.Write(body)
	_ = d.Close()
	err.writeMessageBodyTruncated(w, truncated, clientResponseBodyMaxSizeBytes)
}

func (err *clientResponse) truncateBody(max int) (body []byte, truncated bool) {
	if len(err.value.Body) <= max {
		return err.value.Body, false
	}
	return err.value.Body[:max], true
}

func (err *clientResponse) writeMessageBodyTruncated(w errors.Writer, truncated bool, max int) {
	if truncated {
		_, _ = w.WriteString("(truncated to ")
		_, _ = strconvio.WriteInt(w, int64(max), 10)
		_, _ = w.WriteString(" bytes)\n")
	}
}

func (err *clientResponse) Error() string                 { return errors.Error(err) }
func (err *clientResponse) Format(s fmt.State, verb rune) { errors.Format(err, s, verb) }
func (err *clientResponse) Unwrap() error                 { return err.err }

// GetClientResponse returns the client Response and body associated to an error.
//
// If the error is not wrapped with WithClientResponse(), it returns a nil Response.
// The body may be nil if it was not provided.
func GetClientResponse(err error) *ClientResponse {
	var werr *clientResponse
	ok := errors.As(err, &werr)
	if ok {
		return werr.HTTPClientResponse()
	}
	return nil
}

// ClientResponse represents a client response.
//
// This type is a hack to work around the bodyclose linter.
type ClientResponse struct {
	Response *http.Response
	Body     []byte
}
