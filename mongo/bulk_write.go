package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// WriteModel is a wrapper.
type WriteModel = mongo.WriteModel

// InsertOneModel is a wrapper.
type InsertOneModel = mongo.InsertOneModel

// NewInsertOneModel is a wrapper.
func NewInsertOneModel() *InsertOneModel {
	return mongo.NewInsertOneModel()
}

// DeleteOneModel is a wrapper.
type DeleteOneModel = mongo.DeleteOneModel

// NewDeleteOneModel is a wrapper.
func NewDeleteOneModel() *DeleteOneModel {
	return mongo.NewDeleteOneModel()
}

// DeleteManyModel is a wrapper.
type DeleteManyModel = mongo.DeleteManyModel

// NewDeleteManyModel is a wrapper.
func NewDeleteManyModel() *DeleteManyModel {
	return mongo.NewDeleteManyModel()
}

// ReplaceOneModel is a wrapper.
type ReplaceOneModel = mongo.ReplaceOneModel

// NewReplaceOneModel is a wrapper.
func NewReplaceOneModel() *ReplaceOneModel {
	return mongo.NewReplaceOneModel()
}

// UpdateOneModel is a wrapper.
type UpdateOneModel = mongo.UpdateOneModel

// NewUpdateOneModel is a wrapper.
func NewUpdateOneModel() *UpdateOneModel {
	return mongo.NewUpdateOneModel()
}

// UpdateManyModel is a wrapper.
type UpdateManyModel = mongo.UpdateManyModel

// NewUpdateManyModel is a wrapper.
func NewUpdateManyModel() *UpdateManyModel {
	return mongo.NewUpdateManyModel()
}
