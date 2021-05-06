// Package sibhttpua provides utilities for HTTP outgoing request user-agent.
package sibhttpua

import (
	"fmt"
	"net/http"
)

const format = "Siddhant/1.0 (%s %s; +https://siddhant-test.com)"

// Get returns the user-agent.
func Get(appName string, version string) string {
	check(appName, version)
	return fmt.Sprintf(format, appName, version)
}

func check(appName string, version string) {
	if appName == "" {
		panic("appName must not be empty")
	}
	if version == "" {
		panic("version must not be empty")
	}
}

// SetToRequest sets the user-agent to a request.
func SetToRequest(req *http.Request, appName string, version string) {
	ua := Get(appName, version)
	setUAToRequest(req, ua)
}

func setUAToRequest(req *http.Request, ua string) {
	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", ua)
	}
}

type roundTripper struct {
	http.RoundTripper
	ua string
}

// WrapRoundTripper wraps the given RoundTripper and adds the user-agent to HTTP request.
func WrapRoundTripper(rt http.RoundTripper, appName string, version string) http.RoundTripper {
	ua := Get(appName, version)
	return &roundTripper{
		RoundTripper: rt,
		ua:           ua,
	}
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	setUAToRequest(req, rt.ua)
	return rt.RoundTripper.RoundTrip(req)
}

// WrapDefaultTransport initializes http.DefaultTransport.
// It wraps it with WrapRoundTripper.
func WrapDefaultTransport(appName string, version string) {
	http.DefaultTransport = WrapRoundTripper(http.DefaultTransport, appName, version)
}
