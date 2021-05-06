package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// IndexOptionsBuilder is a wrapper.
type IndexOptionsBuilder = mongo.IndexOptionsBuilder //nolint:staticcheck // This is a wrapper.

// NewIndexOptionsBuilder is a wrapper.
func NewIndexOptionsBuilder() *IndexOptionsBuilder {
	return mongo.NewIndexOptionsBuilder() //nolint:staticcheck // This is a wrapper.
}
