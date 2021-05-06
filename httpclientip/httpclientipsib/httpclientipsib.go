// Package httpclientipsib provides a Getter configured for the infrastructure.
package httpclientipsib

import (
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/httpclientip"
)

var trusted = []string{}

// Getter is a Getter configured for the infrastructure.
var Getter *httpclientip.Getter

func init() {
	g, err := httpclientip.NewGetter(trusted)
	if err != nil {
		err = errors.Wrap(err, "new Getter")
		panic(err)
	}
	Getter = g
}
