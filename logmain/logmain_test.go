package logmain

import (
	"io"
	"log"
	"testing"

	"github.com/siddhant2408/golang-libraries/envutils"
)

func TestStart(t *testing.T) {
	log.SetOutput(io.Discard)
	Start("1.0.0", envutils.Testing)
}
