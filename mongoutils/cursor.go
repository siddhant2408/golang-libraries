package mongoutils

import "context"

// Cursor represents a MongoDB cursor.
type Cursor interface {
	Close(context.Context) error
	Decode(interface{}) error
	All(context.Context, interface{}) error
	Err() error
	Next(context.Context) bool
}
