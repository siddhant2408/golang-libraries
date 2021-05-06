package kafkautils

import "github.com/segmentio/kafka-go"

// AddHeader adds a header.
func AddHeader(hs []kafka.Header, key string, value []byte) []kafka.Header {
	return append(hs, kafka.Header{
		Key:   key,
		Value: value,
	})
}

// SetHeader sets a header.
//
// If the header already exists, the value is replaced, otherwise it is added.
func SetHeader(hs []kafka.Header, key string, value []byte) []kafka.Header {
	for i, h := range hs {
		if h.Key == key {
			h.Value = value
			hs[i] = h
			return hs
		}
	}
	return AddHeader(hs, key, value)
}

// GetHeader returns a header value.
//
// If the header is not defined, the "ok" boolean value is false.
//
// If the header is defined multiple times, the first occurrence is returned.
// See GetHeaderAll.
func GetHeader(hs []kafka.Header, key string) (value []byte, ok bool) {
	for _, h := range hs {
		if h.Key == key {
			return h.Value, true
		}
	}
	return nil, false
}

// GetHeaderAll returns all values for a header.
func GetHeaderAll(hs []kafka.Header, key string) (values [][]byte) {
	for _, h := range hs {
		if h.Key == key {
			values = append(values, h.Value)
		}
	}
	return values
}

// DeleteHeader delates a header.
func DeleteHeader(hs []kafka.Header, key string) []kafka.Header {
	var res []kafka.Header
	for _, h := range hs {
		if h.Key != key {
			res = append(res, h)
		}
	}
	return res
}
