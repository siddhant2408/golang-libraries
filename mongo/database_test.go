package mongo_test

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/mongotest"
	"github.com/siddhant2408/golang-libraries/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDatabaseListCollections(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := coll.InsertOne(ctx, bson.M{"foo": "bar"})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	cu, err := db.ListCollections(ctx, bson.M{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	defer cu.Close(ctx) //nolint:errcheck
	count := 0
	for cu.Next(ctx) {
		count++
	}
	if cu.Err() != nil {
		testutils.FatalErr(t, cu.Err())
	}
	if count != 1 {
		t.Fatalf("unexpected collection count: got %d, want %d", count, 1)
	}
}

func TestDatabaseRunCommand(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := coll.InsertOne(ctx, bson.M{"foo": "bar"})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var v struct {
		N int `bson:"n"`
	}
	err = db.RunCommand(
		ctx,
		bson.M{
			"count": "test",
		},
	).Decode(&v)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if v.N != 1 {
		t.Fatalf("unexpected collection count: got %d, want %d", v.N, 1)
	}
}
