package httphandler

import (
	"net/http"
)

// RequestBodyMaxBytes is a Handler that limits the maximum size of a request body.
type RequestBodyMaxBytes struct {
	http.Handler
	MaxBytes int64
}

// ServeHTTP implements http.Handler.
func (h *RequestBodyMaxBytes) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, h.MaxBytes)
	h.Handler.ServeHTTP(w, req)
}
