package bsonutils

import (
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestConvertString(t *testing.T) {
	for _, tc := range []struct {
		name     string
		v        interface{}
		expected string
	}{
		{
			name:     "String",
			v:        "test",
			expected: "test",
		},
		{
			name:     "Nil",
			v:        nil,
			expected: "",
		},
		{
			name:     "OtherInt",
			v:        int(1),
			expected: "1",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := ConvertString(tc.v)
			if s != tc.expected {
				t.Fatalf("unexpected result: got %q, want %q", s, tc.expected)
			}
		})
	}
}

func TestConvertInt(t *testing.T) {
	for _, tc := range []struct {
		name     string
		v        interface{}
		expected int64
		error    bool
	}{
		{
			name:     "Int64",
			v:        int64(1),
			expected: 1,
		},
		{
			name:     "Int",
			v:        int(1),
			expected: 1,
		},
		{
			name:     "Int32",
			v:        int32(1),
			expected: 1,
		},
		{
			name:     "Float64",
			v:        float64(1),
			expected: 1,
		},
		{
			name:     "Float64Overflow",
			v:        float64(1.23),
			expected: 1,
		},
		{
			name:     "BoolTrue",
			v:        true,
			expected: 1,
		},
		{
			name:     "BoolFalse",
			v:        false,
			expected: 0,
		},
		{
			name:     "String",
			v:        "1",
			expected: 1,
		},
		{
			name:     "Nil",
			v:        nil,
			expected: 0,
		},
		{
			name:     "OtherInt8",
			v:        int8(1),
			expected: 1,
		},
		{
			name:  "ErrorParse",
			v:     "invalid",
			error: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			i, err := ConvertInt(tc.v)
			if err != nil {
				if tc.error {
					return
				}
				testutils.FatalErr(t, err)
			}
			if tc.error {
				t.Fatal("no error")
			}
			if i != tc.expected {
				t.Fatalf("unexpected result: got %d, want %d", i, tc.expected)
			}
		})
	}
}

func TestConvertFloat(t *testing.T) {
	for _, tc := range []struct {
		name     string
		v        interface{}
		expected float64
		error    bool
	}{
		{
			name:     "Float64",
			v:        float64(1.23),
			expected: 1.23,
		},
		{
			name:     "Int",
			v:        int(1),
			expected: 1,
		},
		{
			name:     "Int32",
			v:        int32(1),
			expected: 1,
		},
		{
			name:     "Int64",
			v:        int64(1),
			expected: 1,
		},
		{
			name:     "BoolTrue",
			v:        true,
			expected: 1,
		},
		{
			name:     "BoolFalse",
			v:        false,
			expected: 0,
		},
		{
			name:     "String",
			v:        "1.23",
			expected: 1.23,
		},
		{
			name:     "Nil",
			v:        nil,
			expected: 0,
		},
		{
			name:     "OtherInt8",
			v:        int8(1),
			expected: 1,
		},
		{
			name:  "ErrorParse",
			v:     "invalid",
			error: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, err := ConvertFloat(tc.v)
			if err != nil {
				if tc.error {
					return
				}
				testutils.FatalErr(t, err)
			}
			if tc.error {
				t.Fatal("no error")
			}
			if f != tc.expected {
				t.Fatalf("unexpected result: got %f, want %f", f, tc.expected)
			}
		})
	}
}

func TestConvertBool(t *testing.T) {
	for _, tc := range []struct {
		name     string
		v        interface{}
		expected bool
		error    bool
	}{
		{
			name:     "BoolFalse",
			v:        false,
			expected: false,
		},
		{
			name:     "BoolTrue",
			v:        true,
			expected: true,
		},
		{
			name:     "Int0",
			v:        int(0),
			expected: false,
		},
		{
			name:     "Int1",
			v:        int(1),
			expected: true,
		},
		{
			name:     "Int320",
			v:        int32(0),
			expected: false,
		},
		{
			name:     "Int321",
			v:        int32(1),
			expected: true,
		},
		{
			name:     "Int640",
			v:        int64(0),
			expected: false,
		},
		{
			name:     "Int641",
			v:        int64(1),
			expected: true,
		},
		{
			name:     "Float640",
			v:        float64(0),
			expected: false,
		},
		{
			name:     "Float641",
			v:        float64(1),
			expected: true,
		},
		{
			name:     "String",
			v:        "true",
			expected: true,
		},
		{
			name:     "Nil",
			v:        nil,
			expected: false,
		},
		{
			name:     "OtherInt8",
			v:        int8(1),
			expected: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, err := ConvertBool(tc.v)
			if err != nil {
				if tc.error {
					return
				}
				testutils.FatalErr(t, err)
			}
			if tc.error {
				t.Fatal("no error")
			}
			if b != tc.expected {
				t.Fatalf("unexpected result: got %t, want %t", b, tc.expected)
			}
		})
	}
}
