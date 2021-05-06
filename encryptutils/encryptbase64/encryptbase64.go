// Package encryptbase64 provides an encrypter that converts the output to base64 URL safe.
package encryptbase64

import (
	"encoding/base64"

	"github.com/siddhant2408/golang-libraries/encryptutils"
	"github.com/siddhant2408/golang-libraries/errors"
)

// Encrypter encrypts a []byte to a base64 string (and decrypts reciprocally).
type Encrypter struct {
	enc *encryptutils.Encrypter
}

// NewEncrypter returns a new Encrypter.
func NewEncrypter(key []byte) (*Encrypter, error) {
	enc, err := encryptutils.NewEncrypter(key)
	if err != nil {
		return nil, err
	}
	return &Encrypter{
		enc: enc,
	}, nil
}

// Encrypt encrypts a []byte to a base64 string.
func (e *Encrypter) Encrypt(b []byte) (string, error) {
	b, err := e.enc.Encrypt(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Decrypt decrypts a base64 string to a []byte.
func (e *Encrypter) Decrypt(s string) ([]byte, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, errors.Wrap(err, "base64 decode")
	}
	return e.enc.Decrypt(b)
}
