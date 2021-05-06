package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SingleResult is a wrapper.
type SingleResult struct {
	*singleResult
	database   *Database
	runCommand interface{}
	collection *Collection
	filter     interface{}
}

type singleResult = mongo.SingleResult

// Decode is a wrapper.
func (sr *SingleResult) Decode(v interface{}) (err error) {
	defer sr.wrapErrorReturn("decode", &err, nil)
	return sr.singleResult.Decode(v)
}

// DecodeBytes is a wrapper.
func (sr *SingleResult) DecodeBytes() (_ bson.Raw, err error) {
	defer sr.wrapErrorReturn("decode_bytes", &err, nil)
	return sr.singleResult.DecodeBytes()
}

// Err is a wrapper.
func (sr *SingleResult) Err() (err error) {
	defer sr.wrapErrorReturn("err", &err, nil)
	return sr.singleResult.Err()
}

func (sr *SingleResult) wrapErrorReturn(op string, perr *error, werr func(error) error) {
	wrapErrorReturn("single_result."+op, perr, func(err error) error {
		err = sr.wrapErrorValues(err)
		err = wrapErrorOptional(err, werr)
		return err
	})
}

func (sr *SingleResult) wrapErrorValues(err error) error {
	if sr.database != nil {
		err = sr.database.wrapErrorValues(err)
	}
	if sr.runCommand != nil {
		err = wrapErrorValue(err, "run_command", sr.runCommand)
	}
	if sr.collection != nil {
		err = sr.collection.wrapErrorValues(err)
	}
	if sr.filter != nil {
		err = wrapErrorValue(err, "filter", sr.filter)
	}
	return err
}
