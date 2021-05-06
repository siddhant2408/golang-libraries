// Package encryptutils provides an AES-GCM encrypter.
package encryptutils

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	crypto_rand "crypto/rand"
	"encoding/hex"
	"io"
	"sync"

	"github.com/siddhant2408/golang-libraries/errors"
)

// Encrypter encrypt a []byte to a []byte (and decrypts reciprocally) using AES-GCM.
type Encrypter struct {
	aead cipher.AEAD
}

// NewEncrypter creates a new Encrypter.
func NewEncrypter(key []byte) (*Encrypter, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "new AES cipher")
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "new GCM")
	}
	return &Encrypter{
		aead: aead,
	}, nil
}

// Encrypt encrypts a []byte to a []byte.
func (e *Encrypter) Encrypt(b []byte) ([]byte, error) {
	nonceSize := e.aead.NonceSize()
	nonce := make([]byte, nonceSize)
	bufCrandReader := bufCrandReaderPool.Get().(io.Reader) //nolint:errcheck
	if _, err := io.ReadFull(bufCrandReader, nonce); err != nil {
		return nil, errors.Wrap(err, "generate nonce")
	}
	bufCrandReaderPool.Put(bufCrandReader)
	sealed := e.aead.Seal(nil, nonce, b, nil)
	b = make([]byte, 0, len(nonce)+len(sealed))
	b = append(b, nonce...)
	b = append(b, sealed...)
	return b, nil
}

// Decrypt decrypts a []byte to a []byte.
func (e *Encrypter) Decrypt(b []byte) ([]byte, error) {
	nonceSize := e.aead.NonceSize()
	if len(b) < nonceSize {
		return nil, errors.New("not enough data")
	}
	nonce := b[:nonceSize]
	sealed := b[nonceSize:]
	b, err := e.aead.Open(nil, nonce, sealed, nil)
	if err != nil {
		return nil, errors.Wrap(err, "AEAD decrypt")
	}
	return b, nil
}

// MustHexDecodeString decode an hexx string to a []byte.
// It is suited to decode encryption key.
func MustHexDecodeString(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

var bufCrandReaderPool = &sync.Pool{
	New: func() interface{} {
		return bufio.NewReader(crypto_rand.Reader)
	},
}
