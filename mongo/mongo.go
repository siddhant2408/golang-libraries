// Package mongo is a wrapper for go.mongodb.org/mongo-driver/mongo.
//
// It provides tracing and better errors.
package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// BSONAppender is a wrapper.
type BSONAppender = mongo.BSONAppender

// BSONAppenderFunc is a wrapper.
type BSONAppenderFunc = mongo.BSONAppenderFunc

// Dialer is a wrapper.
type Dialer = mongo.Dialer

// Pipeline is a wrapper.
type Pipeline = mongo.Pipeline
