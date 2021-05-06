package mongoutils_test

import (
	"testing"

	"github.com/siddhant2408/golang-libraries/mongotest"
	"github.com/siddhant2408/golang-libraries/mongoutils"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestForceDisconnect(t *testing.T) {
	_ = mongotest.Connect(t)
}

func TestNewForceDisconnect(t *testing.T) {
	clt := mongotest.Connect(t)
	cl := mongoutils.NewForceDisconnect(clt)
	err := cl()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
