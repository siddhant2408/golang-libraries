package mongotest

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestConnect(t *testing.T) {
	ctx := context.Background()
	clt := Connect(t)
	err := clt.Ping(ctx, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestConnect2(t *testing.T) {
	TestConnect(t)
}

func TestCheckAvailable(t *testing.T) {
	CheckAvailable(t)
}

func TestGetDatabase(t *testing.T) {
	ctx := context.Background()
	db := GetDatabase(t)
	_, err := db.Collection("test").InsertOne(ctx, bson.M{"foo": "bar"})
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestCursor(t *testing.T) {
	ctx := context.Background()
	type myType struct {
		Field string
	}
	docs := []interface{}{
		&myType{
			Field: "a",
		},
		&myType{
			Field: "b",
		},
	}
	c := &Cursor{
		Documents: docs,
	}
	defer c.Close(ctx) //nolint:errcheck
	var resDocs []*myType
	err := c.All(ctx, &resDocs)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if len(resDocs) != len(docs) {
		t.Fatalf("unexpected length: got %d, want %d", len(resDocs), len(docs))
	}
	for i := range resDocs {
		testutils.Compare(t, "unexpected document", resDocs[i], docs[i])
	}
}

func TestCursorEmpty(t *testing.T) {
	ctx := context.Background()
	c := &Cursor{}
	if c.Next(ctx) {
		t.Fatal("next")
	}
}

func TestCursorErr(t *testing.T) {
	c := &Cursor{
		Error: errors.New("error"),
	}
	err := c.Err()
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCursorDecodeErrorResultNotValid(t *testing.T) {
	c := &Cursor{}
	err := c.Decode(nil)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCursorDecodeErrorResultNotPointer(t *testing.T) {
	c := &Cursor{}
	err := c.Decode("invalid")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCursorDecodeErrorResultElementNotValid(t *testing.T) {
	c := &Cursor{}
	err := c.Decode((*string)(nil))
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCursorDecodeErrorBounds(t *testing.T) {
	c := &Cursor{}
	var v string
	err := c.Decode(&v)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCursorDecodeErrorElementNotValid(t *testing.T) {
	ctx := context.Background()
	c := &Cursor{
		Documents: []interface{}{
			nil,
		},
	}
	if !c.Next(ctx) {
		t.Fatal("no next")
	}
	var v string
	err := c.Decode(&v)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCursorDecodeErrorNotAssignable(t *testing.T) {
	ctx := context.Background()
	c := &Cursor{
		Documents: []interface{}{
			1,
		},
	}
	if !c.Next(ctx) {
		t.Fatal("no next")
	}
	var v string
	err := c.Decode(&v)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCursorAllErrorResultNotValid(t *testing.T) {
	ctx := context.Background()
	c := &Cursor{}
	err := c.All(ctx, nil)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCursorAllErrorResultNotPointer(t *testing.T) {
	ctx := context.Background()
	c := &Cursor{}
	err := c.All(ctx, "invalid")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCursorAllErrorResultElementNotValid(t *testing.T) {
	ctx := context.Background()
	c := &Cursor{}
	err := c.All(ctx, (*string)(nil))
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCursorAllErrorResultElementNotSlice(t *testing.T) {
	ctx := context.Background()
	c := &Cursor{}
	var res string
	err := c.All(ctx, &res)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCursorAllErrorDecode(t *testing.T) {
	ctx := context.Background()
	c := &Cursor{
		Documents: []interface{}{
			1,
		},
	}
	var res []string
	err := c.All(ctx, &res)
	if err == nil {
		t.Fatal("no error")
	}
}
