// Package sibhttpapp provides utilities for managing the Applications header used by internal applications.
package sibhttpapp

import (
	"net/http"

	"github.com/siddhant2408/golang-libraries/errors"
)

const header = "Application"

// Set sets the Application header to the request.
func Set(req *http.Request, appName string) {
	if req.Header != nil {
		req.Header.Set(header, appName)
	}
}

// Get gets the Application header from a request.
// It returns an error if the header is not set.
func Get(req *http.Request) (string, error) {
	appName := req.Header.Get(header)
	if appName == "" {
		return "", errors.Newf("header %q is not set", header)
	}
	return appName, nil
}
