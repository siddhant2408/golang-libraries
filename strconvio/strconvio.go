// Package strconvio provides IO utilities for strconv.
package strconvio

import (
	"io"
	"strconv"
	"sync"
)

var bufPool = &sync.Pool{
	New: func() interface{} {
		// 128 should be enough for all the values written by this package.
		return make([]byte, 0, 128)
	},
}

// WriteBool writes a bool.
func WriteBool(w io.Writer, b bool) (n int, err error) {
	bufItf := bufPool.Get()
	defer bufPool.Put(bufItf)
	buf := bufItf.([]byte) //nolint:errcheck
	buf = strconv.AppendBool(buf, b)
	return w.Write(buf)
}

// WriteFloat writes a float.
func WriteFloat(w io.Writer, f float64, format byte, prec, bitSize int) (n int, err error) {
	bufItf := bufPool.Get()
	defer bufPool.Put(bufItf)
	buf := bufItf.([]byte) //nolint:errcheck
	buf = strconv.AppendFloat(buf, f, format, prec, bitSize)
	return w.Write(buf)
}

// WriteInt writes an int.
func WriteInt(w io.Writer, i int64, base int) (n int, err error) {
	bufItf := bufPool.Get()
	defer bufPool.Put(bufItf)
	buf := bufItf.([]byte) //nolint:errcheck
	buf = strconv.AppendInt(buf, i, base)
	return w.Write(buf)
}

// WriteUint writes a uint.
func WriteUint(w io.Writer, i uint64, base int) (n int, err error) {
	bufItf := bufPool.Get()
	defer bufPool.Put(bufItf)
	buf := bufItf.([]byte) //nolint:errcheck
	buf = strconv.AppendUint(buf, i, base)
	return w.Write(buf)
}
