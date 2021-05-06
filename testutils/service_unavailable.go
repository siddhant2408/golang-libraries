package testutils

import (
	"os"
	"strconv"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
)

const unavailableSkipDefaultEnvVar = "TESTUTILS_UNAVAILABLE_SKIP"

// HandleUnavailable handles unavailable services in tests.
//
// It skips or fatals the test depending on the value of boolean environment variables.
// It checks in this order the environment variables:
//  - with the name given in the unavailableSkipEnvVar parameter
//  - TESTUTILS_UNAVAILABLE_SKIP (global default behavior)
func HandleUnavailable(tb testing.TB, unavailableSkipEnvVar string, myerr error) {
	tb.Helper()
	skip, err := shouldSkip(unavailableSkipEnvVar)
	if err != nil {
		err = errors.Wrap(err, "should skip")
		err = errors.Wrap(err, "handle unavailable")
		FatalErr(tb, err)
	}
	if skip {
		SkipErr(tb, myerr)
	} else {
		FatalErr(tb, myerr)
	}
}

func shouldSkip(unavailableSkipEnvVar string) (bool, error) {
	skip, ok, err := getBoolEnvVars(unavailableSkipEnvVar, unavailableSkipDefaultEnvVar)
	if err != nil {
		return false, errors.Wrap(err, "get bool env vars")
	}
	if ok {
		return skip, nil
	}
	// Skip by default.
	return true, nil
}

func getBoolEnvVars(names ...string) (value bool, ok bool, err error) {
	for _, name := range names {
		value, ok, err = getBoolEnvVar(name)
		if err != nil {
			return false, false, errors.Wrap(err, name)
		}
		if ok {
			return value, true, nil
		}
	}
	return false, false, nil
}

func getBoolEnvVar(name string) (value bool, ok bool, err error) {
	s, ok := os.LookupEnv(name)
	if !ok {
		return false, false, nil
	}
	value, err = strconv.ParseBool(s)
	if err != nil {
		return false, false, errors.Wrap(err, "parse bool")
	}
	return value, true, nil
}
