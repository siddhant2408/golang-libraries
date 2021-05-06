// Package jsontest provides JSON utilities for test.
package jsontest

import (
	"encoding/json"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

// Marshal is an alias for MarshalFatal.
func Marshal(tb testing.TB, v interface{}) []byte {
	tb.Helper()
	return MarshalFatal(tb, v)
}

// MarshalFatal marshals the value and returns it.
// It calls Fatal if an error occurs.
func MarshalFatal(tb testing.TB, v interface{}) []byte {
	tb.Helper()
	return marshal(tb, v, testutils.FatalErr)
}

// MarshalError marshals the value and returns it.
// It calls Error if an error occurs.
func MarshalError(tb testing.TB, v interface{}) []byte {
	tb.Helper()
	return marshal(tb, v, testutils.ErrorErr)
}

func marshal(tb testing.TB, v interface{}, f func(testing.TB, error)) []byte {
	tb.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		f(tb, err)
	}
	return b
}

// Unmarshal is an alias UnmarshalFatal.
func Unmarshal(tb testing.TB, data []byte, v interface{}) {
	tb.Helper()
	UnmarshalFatal(tb, data, v)
}

// UnmarshalFatal unmarshals bytes to a value.
// It calls Fatal if an error occurs.
func UnmarshalFatal(tb testing.TB, data []byte, v interface{}) {
	tb.Helper()
	unmarshal(tb, data, v, testutils.FatalErr)
}

// UnmarshalError unmarshals bytes to a value.
// It calls Error if an error occurs.
func UnmarshalError(tb testing.TB, data []byte, v interface{}) {
	tb.Helper()
	unmarshal(tb, data, v, testutils.ErrorErr)
}

func unmarshal(tb testing.TB, data []byte, v interface{}, f func(testing.TB, error)) {
	tb.Helper()
	err := json.Unmarshal(data, v)
	if err != nil {
		f(tb, err)
	}
}
