package bsonutils

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/getsentry/raven-go"
	"github.com/siddhant2408/golang-libraries/errorhandle"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/ravenerrors"
)

// ConvertString converts an interface{} to a string.
func ConvertString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case nil:
		return ""
	}
	report("string", fmt.Sprintf("%T", v))
	return fmt.Sprint(v)
}

// ConvertInt converts an interface{} to an int64.
func ConvertInt(v interface{}) (int64, error) {
	switch v := v.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		report("int64", "string")
		i, err := strconv.ParseInt(v, 10, 64)
		return i, errors.Wrap(err, "parse int")
	case nil:
		return 0, nil
	}
	return ConvertInt(fmt.Sprint(v))
}

// ConvertFloat converts an interface{} to a float64.
func ConvertFloat(v interface{}) (float64, error) {
	switch v := v.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		report("float64", "string")
		f, err := strconv.ParseFloat(v, 64)
		return f, errors.Wrap(err, "parse float")
	case nil:
		return 0, nil
	}
	return ConvertFloat(fmt.Sprint(v))
}

// ConvertBool converts an interface{} to a bool.
func ConvertBool(v interface{}) (bool, error) {
	switch v := v.(type) {
	case bool:
		return v, nil
	case int:
		return v != 0, nil
	case int32:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case float64:
		return v != 0, nil
	case string:
		report("bool", "string")
		b, err := strconv.ParseBool(v)
		return b, errors.Wrap(err, "parse bool")
	case nil:
		return false, nil
	}
	return ConvertBool(fmt.Sprint(v))
}

var (
	// ReportCallPeriod is the period of the reports.
	ReportCallPeriod  int64 = 10000
	reportCallCounter int64 = -1
)

func report(dst string, src string) {
	call := atomic.AddInt64(&reportCallCounter, 1)
	if call%ReportCallPeriod != 0 {
		return
	}
	ctx := context.Background()
	err := errors.New("bsonutils report")
	err = errors.WithTag(err, "mongodb.convert.destination", dst)
	err = errors.WithTag(err, "mongodb.convert.source", src)
	err = ravenerrors.WithSeverity(err, raven.WARNING)
	errorhandle.Handle(ctx, err)
}
