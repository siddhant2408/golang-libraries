package httphandler

import "net/http"

// Ping is a Handler that handles the request GET /_ping and returns status code 200.
type Ping struct {
	http.Handler
}

func (h *Ping) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet && req.URL.Path == "/_ping" {
		return
	}
	h.Handler.ServeHTTP(w, req)
}
