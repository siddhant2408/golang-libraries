// Package httptracinggorilla provides tracing for github.com/gorilla/mux.
package httptracinggorilla

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/siddhant2408/golang-libraries/httptracing"
)

// WrapRouter wraps a mux.Router.
func WrapRouter(rr *mux.Router) http.Handler {
	return httptracing.WrapHandler(rr, NewRouterResourceResolver(rr))
}

// NewRouterResourceResolver returns a new ResourceResolver for a Router.
func NewRouterResourceResolver(rr *mux.Router) httptracing.ResourceResolver {
	return func(req *http.Request) string {
		var match mux.RouteMatch
		route := "unknown"
		if rr.Match(req, &match) && match.Route != nil {
			r, err := match.Route.GetPathTemplate()
			if err == nil {
				route = r
			}
		}
		resource := req.Method + " " + route
		return resource
	}
}
