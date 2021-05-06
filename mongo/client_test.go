package mongo_test

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/mongotest"
	"github.com/siddhant2408/golang-libraries/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestClient(t *testing.T) {
	_ = mongotest.Connect(t)
}

func TestClientListDatabaseNames(t *testing.T) {
	ctx := context.Background()
	clt := mongotest.Connect(t)
	dbNames, err := clt.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if len(dbNames) == 0 {
		t.Fatal("no databases")
	}
}

func TestClientListDatabases(t *testing.T) {
	ctx := context.Background()
	clt := mongotest.Connect(t)
	res, err := clt.ListDatabases(ctx, bson.M{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if len(res.Databases) == 0 {
		t.Fatal("no databases")
	}
}
