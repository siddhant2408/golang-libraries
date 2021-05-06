package mongoutils_test

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/mongotest"
	"github.com/siddhant2408/golang-libraries/mongoutils"
	"github.com/siddhant2408/golang-libraries/testutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCollectionFindOneAndReplaceUpsertRetry(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	res := mongoutils.CollectionFindOneAndReplaceUpsertRetry(
		ctx,
		coll,
		bson.M{
			"foo": "bar",
		},
		bson.M{
			"a": "b",
		},
		options.FindOneAndReplace().SetReturnDocument(options.After).SetUpsert(true),
	)
	err := res.Err()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestCollectionFindOneAndUpdateUpsertRetry(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	res := mongoutils.CollectionFindOneAndUpdateUpsertRetry(
		ctx,
		coll,
		bson.M{
			"foo": "bar",
		},
		bson.M{
			"$set": bson.M{
				"a": "b",
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true),
	)
	err := res.Err()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestCollectionReplaceOneUpsertRetry(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := mongoutils.CollectionReplaceOneUpsertRetry(
		ctx,
		coll,
		bson.M{
			"foo": "bar",
		},
		bson.M{
			"a": "b",
		},
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestCollectionUpdateOneUpsertRetry(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := mongoutils.CollectionUpdateOneUpsertRetry(
		ctx,
		coll,
		bson.M{
			"foo": "bar",
		},
		bson.M{
			"$set": bson.M{
				"a": "b",
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestCollectionUpdateManyUpsertRetry(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := mongoutils.CollectionUpdateManyUpsertRetry(
		ctx,
		coll,
		bson.M{
			"foo": "bar",
		},
		bson.M{
			"$set": bson.M{
				"a": "b",
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
