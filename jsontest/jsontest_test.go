package jsontest

import (
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestMarshalFatal(t *testing.T) {
	v1 := map[string]interface{}{
		"foo": "bar",
	}
	b := Marshal(t, v1)
	var v2 interface{}
	Unmarshal(t, b, &v2)
	testutils.CompareFatal(t, "unexpected value", v2, v1)
}

func TestMarshalError(t *testing.T) {
	v1 := map[string]interface{}{
		"foo": "bar",
	}
	b := MarshalError(t, v1)
	var v2 interface{}
	UnmarshalError(t, b, &v2)
	testutils.CompareFatal(t, "unexpected value", v2, v1)
}
