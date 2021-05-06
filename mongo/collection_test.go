package mongo_test

import (
	"context"
	"sort"
	"testing"

	"github.com/siddhant2408/golang-libraries/mongo"
	"github.com/siddhant2408/golang-libraries/mongotest"
	"github.com/siddhant2408/golang-libraries/testutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCollectionBulkWrite(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	inserted := bson.M{
		"_id": primitive.NewObjectID(),
		"foo": "bar",
	}
	_, err := coll.BulkWrite(
		ctx,
		[]mongo.WriteModel{
			&mongo.InsertOneModel{
				Document: inserted,
			},
		},
	)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var doc bson.M
	err = coll.FindOne(ctx, bson.M{}).Decode(&doc)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	testutils.Compare(t, "unexpected document", doc, inserted)
}

func TestCollectionCountDocuments(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := coll.InsertOne(ctx, bson.M{"foo": "bar"})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	count, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if count != 1 {
		t.Fatalf("unexpected count: got %d, want %d", count, 1)
	}
}

func TestCollectionDeleteMany(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := coll.InsertOne(ctx, bson.M{"foo": "bar"})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	_, err = coll.DeleteMany(ctx, bson.M{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	count, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if count != 0 {
		t.Fatalf("unexpected count: got %d, want %d", count, 0)
	}
}

func TestCollectionDeleteOne(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := coll.InsertOne(ctx, bson.M{"foo": "bar"})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	_, err = coll.DeleteOne(ctx, bson.M{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	count, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if count != 0 {
		t.Fatalf("unexpected count: got %d, want %d", count, 0)
	}
}

func TestCollectionDistinct(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := coll.InsertMany(
		ctx,
		[]interface{}{
			bson.M{"foo": "bar"},
			bson.M{"foo": "bar"},
			bson.M{"foo": "baz"},
		})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	vals, err := coll.Distinct(ctx, "foo", bson.M{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	valsStr := make([]string, len(vals))
	for i, v := range vals {
		s, ok := v.(string)
		if !ok {
			t.Fatal("not a string")
		}
		valsStr[i] = s
	}
	sort.Strings(valsStr)
	expected := []string{"bar", "baz"}
	testutils.Compare(t, "unexpected values", valsStr, expected)
}

func TestCollectionDrop(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := coll.InsertOne(ctx, bson.M{"foo": "bar"})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = coll.Drop(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	count, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if count != 0 {
		t.Fatalf("unexpected count: got %d, want %d", count, 0)
	}
}

func TestCollectionEstimatedDocumentCount(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := coll.InsertOne(ctx, bson.M{"foo": "bar"})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	count, err := coll.EstimatedDocumentCount(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if count != 1 {
		t.Fatalf("unexpected count: got %d, want %d", count, 1)
	}
}

func TestCollectionFind(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	inserted := bson.M{
		"_id": primitive.NewObjectID(),
		"foo": "bar",
	}
	_, err := coll.InsertOne(ctx, inserted)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	cu, err := coll.Find(ctx, bson.M{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	defer cu.Close(ctx) //nolint:errcheck
	for cu.Next(ctx) {
		var doc bson.M
		err := cu.Decode(&doc)
		if err != nil {
			testutils.FatalErr(t, err)
		}
		testutils.Compare(t, "unexpected document", doc, inserted)
	}
	if cu.Err() != nil {
		testutils.FatalErr(t, cu.Err())
	}
}

func TestCollectionFindOne(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	inserted := bson.M{
		"_id": primitive.NewObjectID(),
		"foo": "bar",
	}
	_, err := coll.InsertOne(ctx, inserted)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var doc bson.M
	err = coll.FindOne(ctx, bson.M{}).Decode(&doc)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	testutils.Compare(t, "unexpected document", doc, inserted)
}

func TestCollectionFindOneErrorNotFound(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	var doc bson.M
	err := coll.FindOne(ctx, bson.M{}).Decode(&doc)
	if err == nil {
		t.Fatal("no error")
	}
	if !mongo.IsErrNoDocuments(err) {
		testutils.FatalErr(t, err)
	}
}

func TestCollectionFindOneAndDelete(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	inserted := bson.M{
		"_id": primitive.NewObjectID(),
		"foo": "bar",
	}
	_, err := coll.InsertOne(ctx, inserted)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var doc bson.M
	err = coll.FindOneAndDelete(ctx, bson.M{}).Decode(&doc)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	testutils.Compare(t, "unexpected document", doc, inserted)
	count, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if count != 0 {
		t.Fatalf("unexpected count: got %d, want %d", count, 0)
	}
}

func TestCollectionFindOneAndReplaceBefore(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	id := primitive.NewObjectID()
	inserted := bson.M{
		"_id": id,
		"foo": "bar",
	}
	_, err := coll.InsertOne(ctx, inserted)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	replaced := bson.M{
		"_id": id,
		"foo": "baz",
	}
	var doc bson.M
	err = coll.
		FindOneAndReplace(
			ctx,
			bson.M{},
			replaced,
			options.
				FindOneAndReplace().
				SetReturnDocument(options.Before),
		).
		Decode(&doc)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	testutils.Compare(t, "unexpected document", doc, inserted)
}

func TestCollectionFindOneAndReplaceAfter(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	id := primitive.NewObjectID()
	inserted := bson.M{
		"_id": id,
		"foo": "bar",
	}
	_, err := coll.InsertOne(ctx, inserted)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	replaced := bson.M{
		"_id": id,
		"foo": "baz",
	}
	var doc bson.M
	err = coll.
		FindOneAndReplace(
			ctx,
			bson.M{},
			replaced,
			options.
				FindOneAndReplace().
				SetReturnDocument(options.After),
		).
		Decode(&doc)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	testutils.Compare(t, "unexpected document", doc, replaced)
}

func TestCollectionFindOneAndUpdateBefore(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	id := primitive.NewObjectID()
	inserted := bson.M{
		"_id": id,
		"foo": "bar",
	}
	_, err := coll.InsertOne(ctx, inserted)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var doc bson.M
	err = coll.
		FindOneAndUpdate(
			ctx,
			bson.M{},
			bson.M{
				"$set": bson.M{
					"foo": "baz",
				},
			},
			options.
				FindOneAndUpdate().
				SetReturnDocument(options.Before),
		).
		Decode(&doc)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	testutils.Compare(t, "unexpected document", doc, inserted)
}

func TestCollectionFindOneAndUpdateAfter(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	id := primitive.NewObjectID()
	inserted := bson.M{
		"_id": id,
		"foo": "bar",
	}
	_, err := coll.InsertOne(ctx, inserted)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var doc bson.M
	err = coll.
		FindOneAndUpdate(
			ctx,
			bson.M{},
			bson.M{
				"$set": bson.M{
					"foo": "baz",
				},
			},
			options.
				FindOneAndUpdate().
				SetReturnDocument(options.After),
		).
		Decode(&doc)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	expected := bson.M{
		"_id": id,
		"foo": "baz",
	}
	testutils.Compare(t, "unexpected document", doc, expected)
}

func TestCollectionReplaceOne(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	id := primitive.NewObjectID()
	inserted := bson.M{
		"_id": id,
		"foo": "bar",
	}
	_, err := coll.InsertOne(ctx, inserted)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	replaced := bson.M{
		"_id": id,
		"foo": "baz",
	}
	_, err = coll.ReplaceOne(ctx, bson.M{}, replaced)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var doc bson.M
	err = coll.FindOne(ctx, bson.M{}).Decode(&doc)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	testutils.Compare(t, "unexpected document", doc, replaced)
}

func TestCollectionUpdateMany(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	id := primitive.NewObjectID()
	inserted := bson.M{
		"_id": id,
		"foo": "bar",
	}
	_, err := coll.InsertOne(ctx, inserted)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	_, err = coll.UpdateMany(
		ctx,
		bson.M{},
		bson.M{
			"$set": bson.M{
				"foo": "baz",
			},
		},
	)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var doc bson.M
	err = coll.FindOne(ctx, bson.M{}).Decode(&doc)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	expected := bson.M{
		"_id": id,
		"foo": "baz",
	}
	testutils.Compare(t, "unexpected document", doc, expected)
}

func TestCollectionUpdateOne(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	id := primitive.NewObjectID()
	inserted := bson.M{
		"_id": id,
		"foo": "bar",
	}
	_, err := coll.InsertOne(ctx, inserted)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	_, err = coll.UpdateOne(
		ctx,
		bson.M{},
		bson.M{
			"$set": bson.M{
				"foo": "baz",
			},
		},
	)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var doc bson.M
	err = coll.FindOne(ctx, bson.M{}).Decode(&doc)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	expected := bson.M{
		"_id": id,
		"foo": "baz",
	}
	testutils.Compare(t, "unexpected document", doc, expected)
}
