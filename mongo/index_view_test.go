package mongo_test

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/mongotest"
	"github.com/siddhant2408/golang-libraries/testutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestIndexViewMany(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	iv := coll.Indexes()
	_, err := iv.CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys: bson.D{
					bson.E{Key: "foo", Value: 1},
				},
			},
			{
				Keys: bson.D{
					bson.E{Key: "bar", Value: 1},
				},
			},
		},
	)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	cu, err := iv.List(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	defer cu.Close(ctx) //nolint:errcheck
	type index struct {
		Name string `bson:"name"`
		Key  bson.D `bson:"key"`
	}
	expected := []index{
		{
			Name: "_id_",
			Key: bson.D{
				bson.E{Key: "_id", Value: int32(1)},
			},
		},
		{
			Name: "foo_1",
			Key: bson.D{
				bson.E{Key: "foo", Value: int32(1)},
			},
		},
		{
			Name: "bar_1",
			Key: bson.D{
				bson.E{Key: "bar", Value: int32(1)},
			},
		},
	}
	i := 0
	for cu.Next(ctx) {
		if i >= len(expected) {
			t.Fatalf("unexpected index count %d", i)
		}
		var idx index
		err = cu.Decode(&idx)
		if err != nil {
			testutils.FatalErr(t, err)
		}
		testutils.Compare(t, "unexpected index", idx, expected[i])
		i++
	}
	if cu.Err() != nil {
		testutils.FatalErr(t, cu.Err())
	}
	if i != len(expected) {
		t.Fatalf("unexpected index count: got %d, want %d", i, len(expected))
	}
	_, err = iv.DropAll(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestIndexViewOne(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	iv := coll.Indexes()
	_, err := iv.CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{
				bson.E{Key: "foo", Value: 1},
			},
			Options: options.Index(),
		},
	)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	_, err = iv.DropOne(ctx, "foo_1")
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
