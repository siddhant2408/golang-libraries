package httphandler

import (
	"net/http"
	"sync"
)

// MaxConcurrent is a handler that limits the maximum number of concurrent requests.
type MaxConcurrent struct {
	http.Handler
	Limit               int
	LimitReachedHandler http.Handler // Called if the limit is reached (optional).

	mu      sync.Mutex
	counter int
}

func (h *MaxConcurrent) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !h.acquire() {
		h.handleLimitReached(w, req)
		return
	}
	defer h.release()
	h.Handler.ServeHTTP(w, req)
}

func (h *MaxConcurrent) acquire() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.counter >= h.Limit {
		return false
	}
	h.counter++
	return true
}

func (h *MaxConcurrent) release() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.counter--
}

func (h *MaxConcurrent) handleLimitReached(w http.ResponseWriter, req *http.Request) {
	if h.LimitReachedHandler != nil {
		h.LimitReachedHandler.ServeHTTP(w, req)
		return
	}
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}
