package mgobsonutils

import (
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestParseObjectIdHex(t *testing.T) {
	id1 := bson.NewObjectId()
	s := id1.Hex()
	id2, err := ParseObjectIdHex(s)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if id2 != id1 {
		t.Fatalf("unexpected ObjectId: got %s, want %s", id2, id1)
	}
}

func TestParseObjectIdHexError(t *testing.T) {
	_, err := ParseObjectIdHex("invalid")
	if err == nil {
		t.Fatal("no error")
	}
}
