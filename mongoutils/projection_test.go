package mongoutils_test

import (
	"testing"

	"github.com/siddhant2408/golang-libraries/mongoutils"
	"github.com/siddhant2408/golang-libraries/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestProjection(t *testing.T) {
	expected := bson.M{
		"normal":     1,
		"other_name": 1,
		"noname":     1,
		"notag":      1,
	}
	testProjection(t, func(name string, v interface{}) {
		t.Run(name, func(t *testing.T) {
			sl := mongoutils.Projection(v)
			testutils.Compare(t, "unexpected result", sl, expected)
			sl = mongoutils.Projection(v) // test the cache
			testutils.Compare(t, "unexpected result", sl, expected)
		})
	})
}

func TestProjectionNoCache(t *testing.T) {
	expected := bson.M{
		"normal":     1,
		"other_name": 1,
		"noname":     1,
		"notag":      1,
	}
	testProjection(t, func(name string, v interface{}) {
		t.Run(name, func(t *testing.T) {
			sl := mongoutils.ProjectionNoCache(v)
			testutils.Compare(t, "unexpected result", sl, expected)
		})
	})
}

var benchmarkProjectionResult interface{}

func BenchmarkProjection(b *testing.B) {
	testProjection(b, func(name string, v interface{}) {
		b.Run(name, func(b *testing.B) {
			var res interface{}
			for i := 0; i < b.N; i++ {
				res = mongoutils.Projection(v)
			}
			benchmarkProjectionResult = res
		})
	})
}

func BenchmarkProjectionNoCache(b *testing.B) {
	testProjection(b, func(name string, v interface{}) {
		b.Run(name, func(b *testing.B) {
			var res bson.M
			for i := 0; i < b.N; i++ {
				res = mongoutils.ProjectionNoCache(v)
			}
			benchmarkProjectionResult = res
		})
	})
}

func testProjection(tb testing.TB, f func(name string, v interface{})) {
	tb.Helper()
	type myType struct {
		Normal        string `bson:"normal"`
		DifferentName string `bson:"other_name"`
		NoName        string `bson:",omitempty"`
		NoTag         string
		unexported    string
	}
	_ = myType{}.unexported // Prevents linters to complain about unused field.
	for _, tc := range []struct {
		name string
		v    interface{}
	}{
		{
			name: "Struct",
			v:    myType{},
		},
		{
			name: "Pointer",
			v:    &myType{},
		},
		{
			name: "Slice",
			v:    []myType{},
		},
		{
			name: "Array",
			v:    [3]myType{},
		},
		{
			name: "SlicePointer",
			v:    []*myType{},
		},
	} {
		f(tc.name, tc.v)
	}
}

func TestProjectionPanicUnsupportedType(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("no panic")
		}
	}()
	mongoutils.Projection("invalid")
}
