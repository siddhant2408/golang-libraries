package mongoutils

import (
	"context"

	"github.com/siddhant2408/golang-libraries/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// See https://docs.mongodb.com/manual/reference/method/db.collection.update/
const upsertRetryAttemptsMax = 2

// CollectionFindOneAndReplaceUpsertRetry is a replacement for Collection.FindOneAndReplace with upsert.
// It retries if the upsert fails due to a duplicate key.
func CollectionFindOneAndReplaceUpsertRetry(ctx context.Context, coll *mongo.Collection, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) (res *mongo.SingleResult) {
	upsertRetry(func() error {
		res = coll.FindOneAndReplace(ctx, filter, replacement, opts...)
		return res.Err()
	})
	return res
}

// CollectionFindOneAndUpdateUpsertRetry is a replacement for Collection.FindOneAndUpdate with upsert.
// It retries if the upsert fails due to a duplicate key.
func CollectionFindOneAndUpdateUpsertRetry(ctx context.Context, coll *mongo.Collection, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (res *mongo.SingleResult) {
	upsertRetry(func() error {
		res = coll.FindOneAndUpdate(ctx, filter, update, opts...)
		return res.Err()
	})
	return res
}

// CollectionReplaceOneUpsertRetry is a replacement for Collection.ReplaceOne with upsert.
// It retries if the upsert fails due to a duplicate key.
func CollectionReplaceOneUpsertRetry(ctx context.Context, coll *mongo.Collection, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (res *mongo.UpdateResult, err error) {
	upsertRetry(func() error {
		res, err = coll.ReplaceOne(ctx, filter, replacement, opts...)
		return err
	})
	return res, err
}

// CollectionUpdateOneUpsertRetry is a replacement for Collection.UpdateOne with upsert.
// It retries if the upsert fails due to a duplicate key.
func CollectionUpdateOneUpsertRetry(ctx context.Context, coll *mongo.Collection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error) {
	upsertRetry(func() error {
		res, err = coll.UpdateOne(ctx, filter, update, opts...)
		return err
	})
	return res, err
}

// CollectionUpdateManyUpsertRetry is a replacement for Collection.UpdateMany with upsert.
// It retries if the upsert fails due to a duplicate key.
func CollectionUpdateManyUpsertRetry(ctx context.Context, coll *mongo.Collection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error) {
	upsertRetry(func() error {
		res, err = coll.UpdateMany(ctx, filter, update, opts...)
		return err
	})
	return res, err
}

func upsertRetry(f func() error) {
	for attempts := 1; ; attempts++ {
		err := f()
		if err == nil || !shouldRetryUpsert(attempts, err) {
			return
		}
	}
}

func shouldRetryUpsert(attempts int, err error) bool {
	return attempts < upsertRetryAttemptsMax && isErrorDuplicateKey(err)
}

func isErrorDuplicateKey(err error) bool {
	return mongo.IsWriteErrorCodes(err, 11000)
}
