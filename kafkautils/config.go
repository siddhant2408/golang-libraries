package kafkautils

import (
	"crypto/tls"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
)

// Config represents a common configuration.
//
// It allows to create different types of objects from a common configuration.
// Warning: the created objects are not fully configured.
type Config struct {
	Brokers []string
	TLS     *tls.Config
	SASL    sasl.Mechanism
}

// NewWriter creates a new writer.
//
// The Transport field is always defined with the result of c.NewTransport().
func (c *Config) NewWriter() *kafka.Writer {
	return &kafka.Writer{
		Addr:      kafka.TCP(c.Brokers...),
		Transport: c.NewTransport(),
	}
}

// NewTransport creates a new kafka.Transport.
func (c *Config) NewTransport() *kafka.Transport {
	return &kafka.Transport{
		TLS:  c.TLS,
		SASL: c.SASL,
	}
}

// NewReaderConfig returns a new kafka.ReaderConfig.
//
// The Dialer field is always defined with the result of c.NewDialer().
func (c *Config) NewReaderConfig() kafka.ReaderConfig {
	return kafka.ReaderConfig{
		Brokers: c.Brokers,
		Dialer:  c.NewDialer(),
	}
}

// NewDialer creates a new kafka.Dialer.
//
// It uses kafka.DefaultDialer as model.
func (c *Config) NewDialer() *kafka.Dialer {
	tmp := *kafka.DefaultDialer
	d := &tmp
	d.TLS = c.TLS
	d.SASLMechanism = c.SASL
	return d
}
