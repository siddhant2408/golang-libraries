package mgobsonutils

import (
	"github.com/globalsign/mgo/bson"
	"github.com/siddhant2408/golang-libraries/errors"
)

// ParseObjectIdHex returns an ObjectId from the provided hex representation.
//
// In contrast to ObjectIdHex, it returns an error instead of panic.
func ParseObjectIdHex(s string) (bson.ObjectId, error) { //nolint:golint // Follow the same naming as mgo.
	if !bson.IsObjectIdHex(s) {
		return "", errors.Newf("invalid ObjectId hex %q", s)
	}
	return bson.ObjectIdHex(s), nil
}
