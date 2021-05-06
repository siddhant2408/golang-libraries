package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ClientEncryption is a wrapper.
type ClientEncryption = mongo.ClientEncryption

// NewClientEncryption is a wrapper.
func NewClientEncryption(keyVaultClient *Client, opts ...*options.ClientEncryptionOptions) (*ClientEncryption, error) {
	return mongo.NewClientEncryption(keyVaultClient.client, opts...)
}
