package httphandler

import (
	"net/http"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/panichandle"
)

// Panic is a Handler that captures panic, and call panichandle.
type Panic struct {
	http.Handler
}

func (h *Panic) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Borrowed from net/http.
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		err, ok := r.(error)
		if ok && errors.Is(err, http.ErrAbortHandler) {
			return
		}
		panichandle.Handler(r)
	}()
	h.Handler.ServeHTTP(w, req)
}
