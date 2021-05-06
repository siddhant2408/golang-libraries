package httpclientip

import (
	"net"
	"net/http"
)

// GetFromRequest returns the client IP from a Request.
//
// It is required to use the Handler from this package, otherwise it returns an error.
func GetFromRequest(req *http.Request) (net.IP, error) {
	return GetFromContext(req.Context())
}
