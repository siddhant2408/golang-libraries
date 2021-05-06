// Package strconvbytes provides strconv function that accept bytes instead of string.
//
// The goal is to provide functions that don't require to allocate memory.
// It can be used by applications that need to parse a lot of bytes efficiently.
// Most applications should use strconv.
// Internally this package uses unsafe.
package strconvbytes

import (
	"strconv"
	"unsafe" //nolint:depguard // unsafe is used in order to convert bytes to string efficiently.
)

// ParseBool is a wrapper for strconv.ParseBool.
func ParseBool(b []byte) (bool, error) {
	s := bytesToString(b)
	return strconv.ParseBool(s)
}

// ParseComplex is a wrapper for strconv.ParseComplex.
func ParseComplex(b []byte, bitSize int) (complex128, error) {
	s := bytesToString(b)
	return strconv.ParseComplex(s, bitSize)
}

// ParseFloat is a wrapper for strconv.ParseFloat.
func ParseFloat(b []byte, bitSize int) (float64, error) {
	s := bytesToString(b)
	return strconv.ParseFloat(s, bitSize)
}

// ParseInt is a wrapper for strconv.ParseInt.
func ParseInt(b []byte, base int, bitSize int) (int64, error) {
	s := bytesToString(b)
	return strconv.ParseInt(s, base, bitSize)
}

// ParseUint is a wrapper for strconv.ParseUint.
func ParseUint(b []byte, base int, bitSize int) (uint64, error) {
	s := bytesToString(b)
	return strconv.ParseUint(s, base, bitSize)
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
