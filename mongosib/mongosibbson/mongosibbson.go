// Package mongosibbson provides MongoDB BSON utilities.
package mongosibbson

import (
	"context"
	"reflect" //nolint:depguard // bsoncodec requires to use reflect.
	"strconv"
	"sync/atomic"

	"github.com/getsentry/raven-go"
	"github.com/siddhant2408/golang-libraries/errorhandle"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/ravenerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// NewRegistry returns a new registry that decodes to/from string.
//
// Databases are written by many applications with weakly typed programming languages.
// So it's very common to have fields defined with different types.
//
// It supports:
//  - string => boolean
//  - string => int
//  - string => uint
//  - string => float
//  - string => bool
//  - double => string
//  - int32 => string
//  - int64 => string
//  - boolean => string
//
// For other types, it delegates to the default decoder.
func NewRegistry() *bsoncodec.Registry {
	rb := bson.NewRegistryBuilder()
	Register(rb)
	return rb.Build()
}

// Register registers to a RegistryBuilder.
func Register(rb *bsoncodec.RegistryBuilder) {
	rb.
		RegisterDefaultDecoder(reflect.Bool, bsoncodec.ValueDecoderFunc(decodeBoolean)).
		RegisterDefaultDecoder(reflect.Int, bsoncodec.ValueDecoderFunc(decodeInt)).
		RegisterDefaultDecoder(reflect.Int8, bsoncodec.ValueDecoderFunc(decodeInt)).
		RegisterDefaultDecoder(reflect.Int16, bsoncodec.ValueDecoderFunc(decodeInt)).
		RegisterDefaultDecoder(reflect.Int32, bsoncodec.ValueDecoderFunc(decodeInt)).
		RegisterDefaultDecoder(reflect.Int64, bsoncodec.ValueDecoderFunc(decodeInt)).
		RegisterDefaultDecoder(reflect.Uint, bsoncodec.ValueDecoderFunc(decodeUint)).
		RegisterDefaultDecoder(reflect.Uint8, bsoncodec.ValueDecoderFunc(decodeUint)).
		RegisterDefaultDecoder(reflect.Uint16, bsoncodec.ValueDecoderFunc(decodeUint)).
		RegisterDefaultDecoder(reflect.Uint32, bsoncodec.ValueDecoderFunc(decodeUint)).
		RegisterDefaultDecoder(reflect.Uint64, bsoncodec.ValueDecoderFunc(decodeUint)).
		RegisterDefaultDecoder(reflect.Float32, bsoncodec.ValueDecoderFunc(decodeFloat)).
		RegisterDefaultDecoder(reflect.Float64, bsoncodec.ValueDecoderFunc(decodeFloat)).
		RegisterDefaultDecoder(reflect.String, bsoncodec.ValueDecoderFunc(decodeString))
}

var (
	defaultValueDecoders = bsoncodec.DefaultValueDecoders{}
	uintCodec            = bsoncodec.NewUIntCodec()
	stringCodec          = bsoncodec.NewStringCodec()
)

func decodeBoolean(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if vr.Type() == bsontype.String {
		return decodeBooleanFromString(vr, val)
	}
	return defaultValueDecoders.BooleanDecodeValue(dctx, vr, val)
}

func decodeBooleanFromString(vr bsonrw.ValueReader, val reflect.Value) error {
	reportTypeConversion(val, bsontype.String)
	s, err := vr.ReadString()
	if err != nil {
		return errors.Wrap(err, "read string")
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		reportParseError(err)
		return nil
	}
	val.SetBool(b)
	return nil
}

func decodeInt(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if vr.Type() == bsontype.String {
		return decodeIntFromString(vr, val)
	}
	return defaultValueDecoders.IntDecodeValue(dctx, vr, val)
}

func decodeIntFromString(vr bsonrw.ValueReader, val reflect.Value) error {
	reportTypeConversion(val, bsontype.String)
	s, err := vr.ReadString()
	if err != nil {
		return errors.Wrap(err, "read string")
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		reportParseError(err)
		return nil
	}
	val.SetInt(i)
	return nil
}

func decodeUint(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if vr.Type() == bsontype.String {
		return decodeUintFromString(vr, val)
	}
	return uintCodec.DecodeValue(dctx, vr, val)
}

func decodeUintFromString(vr bsonrw.ValueReader, val reflect.Value) error {
	reportTypeConversion(val, bsontype.String)
	s, err := vr.ReadString()
	if err != nil {
		return errors.Wrap(err, "read string")
	}
	ui, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		reportParseError(err)
		return nil
	}
	val.SetUint(ui)
	return nil
}

func decodeFloat(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if vr.Type() == bsontype.String {
		return decodeFloatFromString(vr, val)
	}
	return defaultValueDecoders.FloatDecodeValue(dctx, vr, val)
}

func decodeFloatFromString(vr bsonrw.ValueReader, val reflect.Value) error {
	reportTypeConversion(val, bsontype.String)
	s, err := vr.ReadString()
	if err != nil {
		return errors.Wrap(err, "read string")
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		reportParseError(err)
		return nil
	}
	val.SetFloat(f)
	return nil
}

func decodeString(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	switch vr.Type() {
	case bsontype.Double:
		return decodeStringFromDouble(vr, val)
	case bsontype.Int32:
		return decodeStringFromInt32(vr, val)
	case bsontype.Int64:
		return decodeStringFromInt64(vr, val)
	case bsontype.Boolean:
		return decodeStringFromBoolean(vr, val)
	}
	return stringCodec.DecodeValue(dctx, vr, val)
}

func decodeStringFromDouble(vr bsonrw.ValueReader, val reflect.Value) error {
	reportTypeConversion(val, bsontype.Double)
	f, err := vr.ReadDouble()
	if err != nil {
		return errors.Wrap(err, "read double")
	}
	s := strconv.FormatFloat(f, 'f', -1, 64)
	val.SetString(s)
	return nil
}

func decodeStringFromInt32(vr bsonrw.ValueReader, val reflect.Value) error {
	reportTypeConversion(val, bsontype.Int32)
	i, err := vr.ReadInt32()
	if err != nil {
		return errors.Wrap(err, "read int32")
	}
	s := strconv.FormatInt(int64(i), 10)
	val.SetString(s)
	return nil
}

func decodeStringFromInt64(vr bsonrw.ValueReader, val reflect.Value) error {
	reportTypeConversion(val, bsontype.Int64)
	i, err := vr.ReadInt64()
	if err != nil {
		return errors.Wrap(err, "read int64")
	}
	s := strconv.FormatInt(i, 10)
	val.SetString(s)
	return nil
}

func decodeStringFromBoolean(vr bsonrw.ValueReader, val reflect.Value) error {
	reportTypeConversion(val, bsontype.Boolean)
	b, err := vr.ReadBoolean()
	if err != nil {
		return errors.Wrap(err, "read boolean")
	}
	s := strconv.FormatBool(b)
	val.SetString(s)
	return nil
}

var (
	// ReportCallPeriod is the period of the reports.
	ReportCallPeriod  int64 = 10000
	reportCallCounter int64 = -1
)

func reportTypeConversion(dst reflect.Value, src bsontype.Type) {
	call := atomic.AddInt64(&reportCallCounter, 1)
	if call%ReportCallPeriod != 0 {
		return
	}
	ctx := context.Background()
	err := errors.New("mongosibbson: type conversion")
	err = errors.WithTag(err, "mongodb.convert.destination", dst.Type().String())
	err = errors.WithTag(err, "mongodb.convert.source", src.String())
	err = ravenerrors.WithSeverity(err, raven.WARNING)
	errorhandle.Handle(ctx, err)
}

func reportParseError(err error) {
	ctx := context.Background()
	err = errors.Wrap(err, "mongosibbson: parse")
	errorhandle.Handle(ctx, err)
}
