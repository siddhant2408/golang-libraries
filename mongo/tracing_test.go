package mongo

import (
	"bytes"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

var marshalJSONTestCases = []struct {
	name     string
	value    interface{}
	expected []byte
}{
	{
		name: "bson.M",
		value: bson.M{
			"foo": "bar",
		},
		expected: []byte(`{"foo":"bar"}`),
	},
	{
		name: "bson.D",
		value: bson.D{
			{
				Key:   "foo",
				Value: "bar",
			},
		},
		expected: []byte(`{"foo":"bar"}`),
	},
	{
		name: "bson.A",
		value: bson.A{
			"test",
		},
		expected: []byte(`["test"]`),
	},
	{
		name: "map[string]string",
		value: map[string]string{
			"foo": "bar",
		},
		expected: []byte(`{"foo":"bar"}`),
	},
	{
		name: "struct",
		value: &struct {
			Foo string `json:"foo"`
		}{
			Foo: "bar",
		},
		expected: []byte(`{"foo":"bar"}`),
	},
	{
		name: "pointer",
		value: &struct {
			Foo string `json:"foo"`
		}{
			Foo: "bar",
		},
		expected: []byte(`{"foo":"bar"}`),
	},
	{
		name:     "string",
		value:    "test",
		expected: []byte(`"test"`),
	},
	{
		name:     "int64",
		value:    int64(123456),
		expected: []byte(`123456`),
	},
	{
		name:     "float64",
		value:    float64(123.456),
		expected: []byte(`123.456`),
	},
	{
		name:     "nil",
		value:    nil,
		expected: []byte(`null`),
	},
}

func TestMarshalJSON(t *testing.T) {
	for _, tc := range marshalJSONTestCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := marshalJSON(tc.value)
			if err != nil {
				testutils.FatalErr(t, err)
			}
			if !bytes.Equal(b, tc.expected) {
				t.Fatalf("unexpected result: got %q, want %q", b, tc.expected)
			}
		})
	}
}

func BenchmarkMarshalJSON(b *testing.B) {
	for _, tc := range marshalJSONTestCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := marshalJSON(tc.value)
				if err != nil {
					testutils.FatalErr(b, err)
				}
			}
		})
	}
}
