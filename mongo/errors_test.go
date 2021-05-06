package mongo_test

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/mongo"
	"github.com/siddhant2408/golang-libraries/mongotest"
	"github.com/siddhant2408/golang-libraries/testutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestErrorTemporaryNoDocuments(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	var doc bson.M
	err := coll.FindOne(ctx, bson.M{}).Decode(&doc)
	if err == nil {
		t.Fatal("no error")
	}
	if errors.IsTemporary(err) {
		t.Fatal("temporary")
	}
}

func TestErrorTemporaryDuplicateKey(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			bson.E{
				Key:   "test",
				Value: 1,
			},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	_, err = coll.InsertOne(ctx, bson.M{"test": "test"})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	_, err = coll.InsertOne(ctx, bson.M{"test": "test"})
	if err == nil {
		t.Fatal("no error")
	}
	if errors.IsTemporary(err) {
		t.Fatal("temporary")
	}
}

func TestErrorTemporaryBadValue(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := coll.InsertOne(ctx, bson.M{"test": "test"})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	_, err = coll.UpdateOne(ctx, bson.M{}, bson.M{
		"$addToSet": bson.M{
			"test": "test",
		},
	})
	if err == nil {
		t.Fatal("no error")
	}
	if errors.IsTemporary(err) {
		t.Fatal("temporary")
	}
}

func TestErrorTemporaryBulk(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			bson.E{
				Key:   "test",
				Value: 1,
			},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	_, err = coll.BulkWrite(ctx, []mongo.WriteModel{
		mongo.NewInsertOneModel().SetDocument(bson.M{"test": "test"}),
		mongo.NewInsertOneModel().SetDocument(bson.M{"test": "test"}),
	})
	if err == nil {
		t.Fatal("no error")
	}
	if errors.IsTemporary(err) {
		t.Fatal("temporary")
	}
}

func TestErrorTemporaryMarshal(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := coll.InsertOne(ctx, bson.M{"test": func() {}})
	if err == nil {
		t.Fatal("no error")
	}
	if errors.IsTemporary(err) {
		t.Fatal("temporary")
	}
}

func TestErrorTemporarySingleResultDecode(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := coll.InsertOne(ctx, bson.M{"test": "test"})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var doc struct {
		Test []string `bson:"test"`
	}
	err = coll.FindOne(ctx, bson.M{}).Decode(&doc)
	if err == nil {
		t.Fatal("no error")
	}
	if errors.IsTemporary(err) {
		t.Fatal("temporary")
	}
}

func TestErrorTemporaryCursorAllDecode(t *testing.T) {
	ctx := context.Background()
	db := mongotest.GetDatabase(t)
	coll := db.Collection("test")
	_, err := coll.InsertOne(ctx, bson.M{"test": "test"})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var res []struct {
		Test []string `bson:"test"`
	}
	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = cur.All(ctx, &res)
	if err == nil {
		t.Fatal("no error")
	}
	if errors.IsTemporary(err) {
		t.Fatal("temporary")
	}
}

func TestIsWriteErrorCodesWriteException(t *testing.T) {
	code := 123
	err := mongo.WriteException{
		WriteErrors: mongo.WriteErrors{
			mongo.WriteError{
				Code: code,
			},
		},
	}
	ok := mongo.IsWriteErrorCodes(err, code)
	if !ok {
		t.Fatal("no match")
	}
}

func TestIsWriteErrorCodesWriteExceptionNoMatch(t *testing.T) {
	err := mongo.WriteException{
		WriteErrors: mongo.WriteErrors{
			mongo.WriteError{
				Code: 123,
			},
		},
	}
	ok := mongo.IsWriteErrorCodes(err, 456)
	if ok {
		t.Fatal("match")
	}
}

func TestIsWriteErrorCodesBulkWriteException(t *testing.T) {
	code := 123
	err := mongo.BulkWriteException{
		WriteErrors: []mongo.BulkWriteError{
			{
				WriteError: mongo.WriteError{
					Code: code,
				},
			},
		},
	}
	ok := mongo.IsWriteErrorCodes(err, code)
	if !ok {
		t.Fatal("no match")
	}
}

func TestIsWriteErrorCodesBulkWriteExceptionNoMatch(t *testing.T) {
	err := mongo.BulkWriteException{
		WriteErrors: []mongo.BulkWriteError{
			{
				WriteError: mongo.WriteError{
					Code: 123,
				},
			},
		},
	}
	ok := mongo.IsWriteErrorCodes(err, 456)
	if ok {
		t.Fatal("match")
	}
}

func TestIsWriteErrorCodesUnknownType(t *testing.T) {
	err := errors.New("error")
	ok := mongo.IsWriteErrorCodes(err, 123)
	if ok {
		t.Fatal("match")
	}
}
