package encryptbase64

import (
	"bytes"
	"testing"

	"github.com/siddhant2408/golang-libraries/encryptutils"
	"github.com/siddhant2408/golang-libraries/testutils"
)

const testHexKey = "a8b41e54aba1cd8fd5015ded3ce7d1f8"

var testData = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

func TestEncrypter(t *testing.T) {
	e := newTestEncrypter(t)
	s, err := e.Encrypt(testData)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	b, err := e.Decrypt(s)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if !bytes.Equal(b, testData) {
		t.Fatalf("unexpected result: got %v, want %v", b, testData)
	}
}

func TestNewEncrypterErrorKey(t *testing.T) {
	_, err := NewEncrypter(nil)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestEncrypterDecryptErrorBase64(t *testing.T) {
	e := newTestEncrypter(t)
	_, err := e.Decrypt(" invalid ")
	if err == nil {
		t.Fatal("no error")
	}
}

func BenchmarkEncrypterEncrypt(b *testing.B) {
	e := newTestEncrypter(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := e.Encrypt(testData)
		if err != nil {
			testutils.FatalErr(b, err)
		}
	}
}

func BenchmarkEncrypterDecrypt(b *testing.B) {
	e := newTestEncrypter(b)
	bb, err := e.Encrypt(testData)
	if err != nil {
		testutils.FatalErr(b, err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := e.Decrypt(bb)
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
