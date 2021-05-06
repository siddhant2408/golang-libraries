package mysqltest

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestGetDatabase(t *testing.T) {
	ctx := context.Background()
	db := GetDatabase(t)
	err := db.PingContext(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestCheckAvailable(t *testing.T) {
	CheckAvailable(t)
}
