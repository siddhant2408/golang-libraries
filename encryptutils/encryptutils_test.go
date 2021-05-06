package encryptutils

import (
	"bytes"
	"testing"

	"github.com/siddhant2408/golang-libraries/testutils"
)

const testHexKey = "a8b41e54aba1cd8fd5015ded3ce7d1f8"

var testData = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

func TestEncrypter(t *testing.T) {
	e := newTestEncrypter(t)
	b, err := e.Encrypt(testData)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	b, err = e.Decrypt(b)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if !bytes.Equal(b, testData) {
		t.Fatalf("unexpected result: got %v, want %v", b, testData)
	}
}

func TestNewEncrypterErrorCipherKey(t *testing.T) {
	_, err := NewEncrypter(nil)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestEncrypterDecryptErrorNotEnoughData(t *testing.T) {
	e := newTestEncrypter(t)
	_, err := e.Decrypt(nil)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestEncrypterDecryptErrorAEADOpen(t *testing.T) {
	e := newTestEncrypter(t)
	_, err := e.Decrypt([]byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"))
	if err == nil {
		t.Fatal("no error")
	}
}

func TestMustHexDecodeStringPanic(t *testing.T) {
	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatal("no panic")
		}
	}()
	MustHexDecodeString("_invalid_")
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
	e, err := NewEncrypter(MustHexDecodeString(testHexKey))
	if err != nil {
		testutils.FatalErr(tb, err)
	}
	return e
}
