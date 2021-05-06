package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// BulkWriteResult is a wrapper.
type BulkWriteResult = mongo.BulkWriteResult

// InsertOneResult is a wrapper.
type InsertOneResult = mongo.InsertOneResult

// InsertManyResult is a wrapper.
type InsertManyResult = mongo.InsertManyResult

// DeleteResult is a wrapper.
type DeleteResult = mongo.DeleteResult

// ListDatabasesResult is a wrapper.
type ListDatabasesResult = mongo.ListDatabasesResult

// DatabaseSpecification is a wrapper.
type DatabaseSpecification = mongo.DatabaseSpecification

// UpdateResult is a wrapper.
type UpdateResult = mongo.UpdateResult
