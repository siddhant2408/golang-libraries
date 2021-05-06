package errors

import (
	std_errors "errors" //nolint:depguard // Required for Is() and As().
	"fmt"

	"github.com/siddhant2408/golang-libraries/errors/internal"
)

// New returns a new error with a message and a stack.
func New(msg string) error {
	return newError(msg)
}

// Newf returns a new error with a formatted message and a stack.
func Newf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return newError(msg)
}

func newError(msg string) error {
	err := internal.NewBase(msg)
	err = withStack(err, 3)
	return err
}

// As calls std_errors.As.
func As(err error, target interface{}) bool {
	return std_errors.As(err, target)
}

// Is calls std_errors.Is.
func Is(err, target error) bool {
	return std_errors.Is(err, target)
}

// Unwrap calls std_errors.Unwrap.
func Unwrap(err error) error {
	return std_errors.Unwrap(err)
}

// UnwrapAll unwraps all nested errors, and returns the last one.
func UnwrapAll(err error) error {
	for {
		werr := Unwrap(err)
		if werr == nil {
			return err
		}
		err = werr
	}
}
