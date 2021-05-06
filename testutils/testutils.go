// Package testutils provides test utilities.
package testutils

import (
	"io"
	"log"

	_ "github.com/siddhant2408/golang-libraries/ballast" // Initializes ballast.
	"github.com/siddhant2408/golang-libraries/httptestlocal"
	_ "github.com/siddhant2408/golang-libraries/randutils" // Initializes random seed.
	_ "github.com/siddhant2408/golang-libraries/spewutils" // Initializes spew config.
	"github.com/siddhant2408/golang-libraries/timeutils"
)

func init() {
	log.SetOutput(io.Discard)
	timeutils.InitFixed()
	httptestlocal.WrapDefaultTransport()
}
