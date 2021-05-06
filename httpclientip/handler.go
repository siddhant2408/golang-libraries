package httpclientip

import (
	"net"
	"net/http"
	"sync"
)

// Handler is an http.Handler that adds the client IP to the HTTP request.
type Handler struct {
	http.Handler
	Getter *Getter
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	v := &handlerContextValue{
		getter: h.Getter,
		req:    req,
	}
	ctx := req.Context()
	ctx = setToContext(ctx, v)
	req = req.WithContext(ctx)
	h.Handler.ServeHTTP(w, req)
}

type handlerContextValue struct {
	once sync.Once
	ip   net.IP
	err  error

	getter *Getter
	req    *http.Request
}

func (v *handlerContextValue) getClientIP() (net.IP, error) {
	v.once.Do(v.init)
	return v.ip, v.err
}

func (v *handlerContextValue) init() {
	v.ip, v.err = v.getter.GetClientIP(v.req)
}
