// Package encrypturl provides an encrypter for url.Values.
package encrypturl

import (
	"net/url"

	"github.com/siddhant2408/golang-libraries/encryptutils/encryptbase64"
	"github.com/siddhant2408/golang-libraries/errors"
)

// Encrypter encrypts a url.Values to a base64 string (and decrypts reciprocally).
type Encrypter struct {
	enc *encryptbase64.Encrypter
}

// NewEncrypter returns a new Encrypter.
func NewEncrypter(key []byte) (*Encrypter, error) {
	enc, err := encryptbase64.NewEncrypter(key)
	if err != nil {
		return nil, err
	}
	return &Encrypter{
		enc: enc,
	}, nil
}

// Encrypt encrypts a url.Values to a base64 string.
func (e *Encrypter) Encrypt(u url.Values) (string, error) {
	return e.enc.Encrypt([]byte(u.Encode()))
}

// Decrypt decrypts a base64 string to a url.Values.
func (e *Encrypter) Decrypt(s string) (url.Values, error) {
	b, err := e.enc.Decrypt(s)
	if err != nil {
		return nil, err
	}
	q, err := url.ParseQuery(string(b))
	if err != nil {
		return nil, errors.Wrap(err, "parse URL query")
	}
	return q, nil
}
