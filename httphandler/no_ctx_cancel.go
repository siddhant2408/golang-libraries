package httphandler

import (
	"net/http"

	"github.com/siddhant2408/golang-libraries/ctxutils"
)

// NoContextCancel is a Handler that doesn't close the context of the request (if the connection is closed).
type NoContextCancel struct {
	http.Handler
}

// ServeHTTP implements http.Handler.
func (h *NoContextCancel) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = ctxutils.WithoutCancel(ctx)
	req = req.WithContext(ctx)
	h.Handler.ServeHTTP(w, req)
}
