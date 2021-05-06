package sibvalid

import (
	"testing"
)

var testPhonesValid = []string{
	"330123456789",
}

var testPhonesInvalid = []string{
	"0123456789",
	"1",
	"11111",
	"11111111111111111",
	"invalid",
}

func TestPhoneValid(t *testing.T) {
	testStrings(t, testPhoneValid, testPhonesValid)
}

func BenchmarkPhoneValid(b *testing.B) {
	benchmarkStrings(b, testPhoneValid, testPhonesValid)
}

func testPhoneValid(tb testing.TB, s string) {
	testStringValid(tb, Phone, s)
}

func TestPhoneInvalid(t *testing.T) {
	testStrings(t, testPhoneInvalid, testPhonesInvalid)
}

func BenchmarkPhoneInvalid(b *testing.B) {
	benchmarkStrings(b, testPhoneInvalid, testPhonesInvalid)
}

func testPhoneInvalid(tb testing.TB, s string) {
	testStringInvalid(tb, Phone, s)
}
