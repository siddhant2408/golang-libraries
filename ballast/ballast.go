// Package ballast provides a memory ballast implementation.
//
// See:
//  - https://github.com/golang/go/issues/23044
//  - https://blog.twitch.tv/go-memory-ballast-how-i-learnt-to-stop-worrying-and-love-the-heap-26c2462549a2
//
// Importing this package allocates a []byte of the size of the GO_BALLAST environment variable.
package ballast

import (
	"os"
	"strconv"

	"github.com/siddhant2408/golang-libraries/errors"
)

const envVar = "GO_BALLAST"

var (
	alloc []byte
	_     = alloc // Prevents linters to complain about unused variable.
)

func init() {
	s, ok := os.LookupEnv(envVar)
	if !ok {
		return
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		err = errors.Wrapf(err, "parse %s environment variable", envVar)
		panic(err)
	}
	alloc = make([]byte, i)
}
