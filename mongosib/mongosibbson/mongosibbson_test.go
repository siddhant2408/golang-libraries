package mongosibbson

import (
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

type testType struct {
	Test interface{} `bson:"test"`
}

type testTypeBool struct {
	Test bool `bson:"test"`
}

func TestBoolean(t *testing.T) {
	r := NewRegistry()
	for _, tc := range []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "BooleanTrue",
			value:    true,
			expected: true,
		},
		{
			name:     "BooleanFalse",
			value:    false,
			expected: false,
		},
		{
			name:     "StringTrue",
			value:    "true",
			expected: true,
		},
		{
			name:     "StringFalse",
			value:    "false",
			expected: false,
		},
		{
			name:     "String1",
			value:    "1",
			expected: true,
		},
		{
			name:     "String0",
			value:    "0",
			expected: false,
		},
		{
			name:     "Int32-1",
			value:    int32(1),
			expected: true,
		},
		{
			name:     "Int32-0",
			value:    int32(0),
			expected: false,
		},
		{
			name:     "Int64-1",
			value:    int64(1),
			expected: true,
		},
		{
			name:     "Int64-0",
			value:    int64(0),
			expected: false,
		},
		{
			name:     "Double1",
			value:    float64(1),
			expected: true,
		},
		{
			name:     "Double0",
			value:    float64(0),
			expected: false,
		},
		{
			name:     "Null",
			value:    nil,
			expected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, err := bson.Marshal(testType{
				Test: tc.value,
			})
			if err != nil {
				testutils.FatalErr(t, err)
			}
			var v testTypeBool
			err = bson.UnmarshalWithRegistry(r, b, &v)
			if err != nil {
				testutils.FatalErr(t, err)
			}
			if v.Test != tc.expected {
				t.Fatalf("unexpected result: got %t, want %t", v.Test, tc.expected)
			}
		})
	}
}

func TestBooleanErrorStringParse(t *testing.T) {
	r := NewRegistry()
	b, err := bson.Marshal(testType{
		Test: "invalid",
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var v testTypeBool
	err = bson.UnmarshalWithRegistry(r, b, &v)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if v.Test {
		t.Fatalf("unexpected result: got %t, want %t", v.Test, false)
	}
}

func TestBooleanErrorUnsupported(t *testing.T) {
	r := NewRegistry()
	b, err := bson.Marshal(testType{
		Test: []string{"test"},
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var v testTypeBool
	err = bson.UnmarshalWithRegistry(r, b, &v)
	if err == nil {
		t.Fatal("no error")
	}
}

type testTypeInt struct {
	Test int64 `bson:"test"`
}

func TestInt(t *testing.T) {
	r := NewRegistry()
	for _, tc := range []struct {
		name     string
		value    interface{}
		expected int64
	}{
		{
			name:     "Int64",
			value:    int64(1),
			expected: 1,
		},
		{
			name:     "Int32",
			value:    int32(1),
			expected: 1,
		},
		{
			name:     "Double",
			value:    float64(1),
			expected: 1,
		},
		{
			name:     "String",
			value:    "1",
			expected: 1,
		},
		{
			name:     "BooleanTrue",
			value:    true,
			expected: 1,
		},
		{
			name:     "BooleanFalse",
			value:    false,
			expected: 0,
		},
		{
			name:     "Null",
			value:    nil,
			expected: 0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, err := bson.Marshal(testType{
				Test: tc.value,
			})
			if err != nil {
				testutils.FatalErr(t, err)
			}
			var v testTypeInt
			err = bson.UnmarshalWithRegistry(r, b, &v)
			if err != nil {
				testutils.FatalErr(t, err)
			}
			if v.Test != tc.expected {
				t.Fatalf("unexpected result: got %d, want %d", v.Test, tc.expected)
			}
		})
	}
}

func TestIntErrorStringParse(t *testing.T) {
	r := NewRegistry()
	b, err := bson.Marshal(testType{
		Test: "invalid",
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var v testTypeInt
	err = bson.UnmarshalWithRegistry(r, b, &v)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if v.Test != 0 {
		t.Fatalf("unexpected result: got %d, want %d", v.Test, int64(0))
	}
}

func TestIntErrorUnsupported(t *testing.T) {
	r := NewRegistry()
	b, err := bson.Marshal(testType{
		Test: []string{"test"},
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var v testTypeInt
	err = bson.UnmarshalWithRegistry(r, b, &v)
	if err == nil {
		t.Fatal("no error")
	}
}

type testTypeUint struct {
	Test uint64 `bson:"test"`
}

func TestUint(t *testing.T) {
	r := NewRegistry()
	for _, tc := range []struct {
		name     string
		value    interface{}
		expected uint64
	}{
		{
			name:     "Int64",
			value:    int64(1),
			expected: 1,
		},
		{
			name:     "Int32",
			value:    int32(1),
			expected: 1,
		},
		{
			name:     "Double",
			value:    float64(1),
			expected: 1,
		},
		{
			name:     "String",
			value:    "1",
			expected: 1,
		},
		{
			name:     "BooleanTrue",
			value:    true,
			expected: 1,
		},
		{
			name:     "BooleanFalse",
			value:    false,
			expected: 0,
		},
		{
			name:     "Null",
			value:    nil,
			expected: 0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, err := bson.Marshal(testType{
				Test: tc.value,
			})
			if err != nil {
				testutils.FatalErr(t, err)
			}
			var v testTypeUint
			err = bson.UnmarshalWithRegistry(r, b, &v)
			if err != nil {
				testutils.FatalErr(t, err)
			}
			if v.Test != tc.expected {
				t.Fatalf("unexpected result: got %d, want %d", v.Test, tc.expected)
			}
		})
	}
}

func TestUintErrorStringParse(t *testing.T) {
	r := NewRegistry()
	b, err := bson.Marshal(testType{
		Test: "invalid",
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var v testTypeUint
	err = bson.UnmarshalWithRegistry(r, b, &v)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if v.Test != 0 {
		t.Fatalf("unexpected result: got %d, want %d", v.Test, uint64(0))
	}
}

func TestUintErrorUnsupported(t *testing.T) {
	r := NewRegistry()
	b, err := bson.Marshal(testType{
		Test: []string{"test"},
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var v testTypeUint
	err = bson.UnmarshalWithRegistry(r, b, &v)
	if err == nil {
		t.Fatal("no error")
	}
}

type testTypeFloat struct {
	Test float64 `bson:"test"`
}

func TestFloat(t *testing.T) {
	r := NewRegistry()
	for _, tc := range []struct {
		name     string
		value    interface{}
		expected float64
	}{
		{
			name:     "Double",
			value:    float64(1),
			expected: 1,
		},
		{
			name:     "Int32",
			value:    int32(1),
			expected: 1,
		},
		{
			name:     "Int64",
			value:    int64(1),
			expected: 1,
		},
		{
			name:     "String",
			value:    "1",
			expected: 1,
		},
		{
			name:     "BooleanTrue",
			value:    true,
			expected: 1,
		},
		{
			name:     "BooleanFalse",
			value:    false,
			expected: 0,
		},
		{
			name:     "Null",
			value:    nil,
			expected: 0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, err := bson.Marshal(testType{
				Test: tc.value,
			})
			if err != nil {
				testutils.FatalErr(t, err)
			}
			var v testTypeFloat
			err = bson.UnmarshalWithRegistry(r, b, &v)
			if err != nil {
				testutils.FatalErr(t, err)
			}
			if v.Test != tc.expected {
				t.Fatalf("unexpected result: got %f, want %f", v.Test, tc.expected)
			}
		})
	}
}

func TestFloatErrorStringParse(t *testing.T) {
	r := NewRegistry()
	b, err := bson.Marshal(testType{
		Test: "invalid",
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var v testTypeFloat
	err = bson.UnmarshalWithRegistry(r, b, &v)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if v.Test != 0 {
		t.Fatalf("unexpected result: got %f, want %f", v.Test, float64(0))
	}
}

func TestFloatErrorUnsupported(t *testing.T) {
	r := NewRegistry()
	b, err := bson.Marshal(testType{
		Test: []string{"test"},
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var v testTypeFloat
	err = bson.UnmarshalWithRegistry(r, b, &v)
	if err == nil {
		t.Fatal("no error")
	}
}

type testTypeString struct {
	Test string `bson:"test"`
}

func TestString(t *testing.T) {
	r := NewRegistry()
	for _, tc := range []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "String",
			value:    "test",
			expected: "test",
		},
		{
			name:     "Double",
			value:    float64(1),
			expected: "1",
		},
		{
			name:     "Int32",
			value:    int32(1),
			expected: "1",
		},
		{
			name:     "Int64",
			value:    int64(1),
			expected: "1",
		},
		{
			name:     "BooleanTrue",
			value:    true,
			expected: "true",
		},
		{
			name:     "BooleanFalse",
			value:    false,
			expected: "false",
		},
		{
			name:     "Null",
			value:    nil,
			expected: "",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, err := bson.Marshal(testType{
				Test: tc.value,
			})
			if err != nil {
				testutils.FatalErr(t, err)
			}
			var v testTypeString
			err = bson.UnmarshalWithRegistry(r, b, &v)
			if err != nil {
				testutils.FatalErr(t, err)
			}
			if v.Test != tc.expected {
				t.Fatalf("unexpected result: got %q, want %q", v.Test, tc.expected)
			}
		})
	}
}

func TestStringErrorUnsupported(t *testing.T) {
	r := NewRegistry()
	b, err := bson.Marshal(testType{
		Test: []string{"test"},
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var v testTypeString
	err = bson.UnmarshalWithRegistry(r, b, &v)
	if err == nil {
		t.Fatal("no error")
	}
}
