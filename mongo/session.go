package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// For now, we will not wrap the session.

// WithSession is a wrapper.
func WithSession(ctx context.Context, sess Session, fn func(SessionContext) error) error {
	return mongo.WithSession(ctx, sess, fn)
}

// Session is a wrapper.
type Session = mongo.Session

// SessionContext is a wrapper.
type SessionContext = mongo.SessionContext
