package strconvbytes

import (
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestParseBool(t *testing.T) {
	bs := []byte("true")
	res, err := ParseBool(bs)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	expected := true
	if res != expected {
		t.Fatalf("unexpected result: got %t, want %t", res, expected)
	}
}

func BenchmarkParseBool(b *testing.B) {
	bs := []byte("true")
	for i := 0; i < b.N; i++ {
		_, err := ParseBool(bs)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
}

func TestParseComplex(t *testing.T) {
	bs := []byte("1+1i")
	res, err := ParseComplex(bs, 128)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	expected := complex(1, 1)
	if res != expected {
		t.Fatalf("unexpected result: got %f, want %f", res, expected)
	}
}

func BenchmarkParseComplex(b *testing.B) {
	bs := []byte("1+1i")
	for i := 0; i < b.N; i++ {
		_, err := ParseComplex(bs, 128)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
}

func TestParseFloat(t *testing.T) {
	bs := []byte("123.456")
	res, err := ParseFloat(bs, 64)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	expected := 123.456
	if res != expected {
		t.Fatalf("unexpected result: got %f, want %f", res, expected)
	}
}

func BenchmarkParseFloat(b *testing.B) {
	bs := []byte("123.456")
	for i := 0; i < b.N; i++ {
		_, err := ParseFloat(bs, 64)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
}

func TestParseInt(t *testing.T) {
	bs := []byte("123456")
	res, err := ParseInt(bs, 10, 64)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	expected := int64(123456)
	if res != expected {
		t.Fatalf("unexpected result: got %d, want %d", res, expected)
	}
}

func BenchmarkParseInt(b *testing.B) {
	bs := []byte("123456")
	for i := 0; i < b.N; i++ {
		_, err := ParseInt(bs, 10, 64)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
}

func TestParseUint(t *testing.T) {
	bs := []byte("123456")
	res, err := ParseUint(bs, 10, 64)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	expected := uint64(123456)
	if res != expected {
		t.Fatalf("unexpected result: got %d, want %d", res, expected)
	}
}

func BenchmarkParseUint(b *testing.B) {
	bs := []byte("123456")
	for i := 0; i < b.N; i++ {
		_, err := ParseUint(bs, 10, 64)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
}
