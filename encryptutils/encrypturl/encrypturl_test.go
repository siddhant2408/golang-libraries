package encrypturl

import (
	"net/url"
	"testing"

	"github.com/siddhant2408/golang-libraries/encryptutils"
	"github.com/siddhant2408/golang-libraries/testutils"
)

const testHexKey = "a8b41e54aba1cd8fd5015ded3ce7d1f8"

var testURLValues = url.Values{
	"a": {"aaaaaaaaaaaaaaaa"},
	"b": {"bbbbbbbbbbbbbbbb"},
	"c": {"cccccccccccccccc"},
	"d": {"dddddddddddddddd"},
}

func TestEncrypter(t *testing.T) {
	e := newTestEncrypter(t)
	s, err := e.Encrypt(testURLValues)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	vs, err := e.Decrypt(s)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	testutils.Compare(t, "unexpected URL values", vs, testURLValues)
}

func TestNewEncrypterErrorKey(t *testing.T) {
	_, err := NewEncrypter(nil)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestEncrypterDecryptParent(t *testing.T) {
	e := newTestEncrypter(t)
	_, err := e.Decrypt("")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestEncrypterDecryptErrorURLParseQuery(t *testing.T) {
	e := newTestEncrypter(t)
	// This is "%é" correctly encrypted, which is an invalid query string.
	_, err := e.Decrypt("wh5T2_UumQzISI9gq_5-eUMdSlaj0pDYF2DcRfryBg")
	if err == nil {
		t.Fatal("no error")
	}
}

func BenchmarkEncrypterEncrypt(b *testing.B) {
	e := newTestEncrypter(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := e.Encrypt(testURLValues)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
}

func BenchmarkEncrypterDecrypt(b *testing.B) {
	e := newTestEncrypter(b)
	s, err := e.Encrypt(testURLValues)
	if err != nil {
		testutils.FatalErr(b, err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := e.Decrypt(s)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
}

func newTestEncrypter(tb testing.TB) *Encrypter {
	tb.Helper()
	e, err := NewEncrypter(encryptutils.MustHexDecodeString(testHexKey))
	if err != nil {
		testutils.FatalErr(tb, err)
	}
	return e
}
