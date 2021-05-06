// Package randutils provides random related utilities.
package randutils

import (
	"math/rand"

	"github.com/siddhant2408/golang-libraries/timeutils"
)

// For now, it only initializes the global random seed with the current time.

func init() {
	rand.Seed(timeutils.Now().UnixNano())
}
