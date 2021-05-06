package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// ChangeStream is a wrapper.
type ChangeStream struct {
	*changeStream
	client     *Client
	database   *Database
	collection *Collection
	pipeline   interface{}
}

type changeStream = mongo.ChangeStream

// Close is a wrapper.
func (cs *ChangeStream) Close(ctx context.Context) (err error) {
	defer cs.wrapErrorReturn("close", &err, nil)
	return cs.changeStream.Close(ctx)
}

// Decode is a wrapper.
func (cs *ChangeStream) Decode(out interface{}) (err error) {
	defer cs.wrapErrorReturn("decode", &err, nil)
	return cs.changeStream.Decode(out)
}

// Err is a wrapper.
func (cs *ChangeStream) Err() (err error) {
	defer cs.wrapErrorReturn("err", &err, nil)
	return cs.changeStream.Err()
}

func (cs *ChangeStream) wrapErrorReturn(op string, perr *error, werr func(error) error) {
	wrapErrorReturn("change_stream."+op, perr, func(err error) error {
		err = cs.wrapErrorValues(err)
		err = wrapErrorOptional(err, werr)
		return err
	})
}

func (cs *ChangeStream) wrapErrorValues(err error) error {
	if cs.client != nil {
		err = cs.client.wrapErrorValues(err)
	}
	if cs.database != nil {
		err = cs.database.wrapErrorValues(err)
	}
	if cs.collection != nil {
		err = cs.collection.wrapErrorValues(err)
	}
	if cs.pipeline != nil {
		err = wrapErrorValue(err, "pipeline", cs.pipeline)
	}
	return err
}

// StreamType is a wrapper.
type StreamType = mongo.StreamType

// StreamType constants.
const (
	CollectionStream = mongo.CollectionStream
	DatabaseStream   = mongo.DatabaseStream
	ClientStream     = mongo.ClientStream
)
