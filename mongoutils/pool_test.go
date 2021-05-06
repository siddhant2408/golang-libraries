package mongoutils_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/siddhant2408/golang-libraries/mongotest"
	"github.com/siddhant2408/golang-libraries/mongoutils"
	"github.com/siddhant2408/golang-libraries/testutils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestClientPool(t *testing.T) {
	mongotest.CheckAvailable(t)
	ctx := context.Background()
	cp := mongoutils.NewClientPool()
	defer cp.ForceClose(func(err error) {
		testutils.ErrorErr(t, err)
	})
	clt, err := cp.GetClient(ctx, mongotest.GetURI())
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = clt.Ping(ctx, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	clt, err = cp.GetClient(ctx, mongotest.GetURI())
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = clt.Ping(ctx, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestClientPoolError(t *testing.T) {
	ctx := context.Background()
	cp := mongoutils.NewClientPool()
	defer cp.Close(ctx, func(err error) {
		testutils.ErrorErr(t, err)
	})
	_, err := cp.GetClient(ctx, "invalid")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestDatabasePool(t *testing.T) {
	mongotest.CheckAvailable(t)
	ctx := context.Background()
	cp := mongoutils.NewClientPool()
	defer cp.Close(ctx, func(err error) {
		testutils.ErrorErr(t, err)
	})
	dp := &mongoutils.DatabasePool{
		ClientPool: cp,
	}
	for i := 0; i < 3; i++ { // Get 3 databases.
		pu, err := url.Parse(mongotest.GetURI())
		if err != nil {
			testutils.FatalErr(t, err)
		}
		pu.Path = fmt.Sprintf("test_%s", primitive.NewObjectID().Hex())
		u := pu.String()
		db, err := dp.GetDatabase(ctx, u)
		if err != nil {
			testutils.FatalErr(t, err)
		}
		err = db.Drop(ctx)
		if err != nil {
			testutils.FatalErr(t, err)
		}
	}
}

func TestDatabasePoolErrorParseURI(t *testing.T) {
	ctx := context.Background()
	dp := &mongoutils.DatabasePool{}
	_, err := dp.GetDatabase(ctx, ":invalid")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestDatabasePoolErrorParseURIEmptyDBName(t *testing.T) {
	ctx := context.Background()
	dp := &mongoutils.DatabasePool{}
	_, err := dp.GetDatabase(ctx, mongotest.GetURI())
	if err == nil {
		t.Fatal("no error")
	}
}

func TestDatabasePoolErrorClient(t *testing.T) {
	ctx := context.Background()
	cp := mongoutils.NewClientPool()
	defer cp.Close(ctx, func(err error) {
		testutils.ErrorErr(t, err)
	})
	dp := &mongoutils.DatabasePool{
		ClientPool: cp,
	}
	_, err := dp.GetDatabase(ctx, "invalid://test/test")
	if err == nil {
		t.Fatal("no error")
	}
}
