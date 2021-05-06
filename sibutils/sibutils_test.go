package sibutils

import (
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestEncodeDecode(t *testing.T) {
	v1 := []int64{0, 1, 12345, 11111, 99999}
	ref, err := Encode(v1)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	v2, err := Decode(ref)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	testutils.Compare(t, "unexpected result", v2, v1)
}

func TestEncodeEmpty(t *testing.T) {
	s, err := Encode([]int64{})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if s != "" {
		t.Fatalf("unexpected result: got %q, want %q", s, "")
	}
}

func TestDecodeEmptyField(t *testing.T) {
	vi, err := Decode("1jw72y2b0ahmzh6q1im")
	if err != nil {
		testutils.FatalErr(t, err)
	}
	expected := []int64{1495442598290162, 1, 1304717, 0}
	testutils.Compare(t, "unexpected result", vi, expected)
}

func TestDecodeEmpty(t *testing.T) {
	vi, err := Decode("")
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if len(vi) != 0 {
		t.Fatalf("unexpected result length: got %d, want %d", len(vi), 0)
	}
}

func TestEncodeEmptyErrorNegative(t *testing.T) {
	_, err := Encode([]int64{-1})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestDecodeErrorConvertBase(t *testing.T) {
	ref := "#invalid"
	_, err := Decode(ref)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestDecodeErrorParseInt(t *testing.T) {
	ref := "b"
	_, err := Decode(ref)
	if err == nil {
		t.Fatal("no error")
	}
}
