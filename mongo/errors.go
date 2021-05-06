package mongo

import (
	"github.com/siddhant2408/golang-libraries/errors"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
)

// IsDuplicateKeyError is a wrapper.
func IsDuplicateKeyError(err error) bool {
	return mongo.IsDuplicateKeyError(err)
}

// IsNetworkError is a wrapper.
func IsNetworkError(err error) bool {
	return mongo.IsNetworkError(err)
}

// IsTimeout is a wrapper.
func IsTimeout(err error) bool {
	return mongo.IsTimeout(err)
}

// IsErrClientDisconnected is a wrapper.
func IsErrClientDisconnected(err error) bool {
	return errors.Is(err, mongo.ErrClientDisconnected)
}

// IsErrEmptySlice is a wrapper.
func IsErrEmptySlice(err error) bool {
	return errors.Is(err, mongo.ErrEmptySlice)
}

// IsErrInvalidIndexValue is a wrapper.
func IsErrInvalidIndexValue(err error) bool {
	return errors.Is(err, mongo.ErrInvalidIndexValue)
}

// IsErrMissingResumeToken is a wrapper.
func IsErrMissingResumeToken(err error) bool {
	return errors.Is(err, mongo.ErrMissingResumeToken)
}

// IsErrMultipleIndexDrop is a wrapper.
func IsErrMultipleIndexDrop(err error) bool {
	return errors.Is(err, mongo.ErrMultipleIndexDrop)
}

// IsErrNilCursor is a wrapper.
func IsErrNilCursor(err error) bool {
	return errors.Is(err, mongo.ErrNilCursor)
}

// IsErrNilDocument is a wrapper.
func IsErrNilDocument(err error) bool {
	return errors.Is(err, mongo.ErrNilDocument)
}

// IsErrNilValue is a wrapper.
func IsErrNilValue(err error) bool {
	return errors.Is(err, mongo.ErrNilValue)
}

// IsErrNoDocuments is a wrapper.
func IsErrNoDocuments(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}

// IsErrNonStringIndexName is a wrapper.
func IsErrNonStringIndexName(err error) bool {
	return errors.Is(err, mongo.ErrNonStringIndexName)
}

// IsErrUnacknowledgedWrite is a wrapper.
func IsErrUnacknowledgedWrite(err error) bool {
	return errors.Is(err, mongo.ErrUnacknowledgedWrite)
}

// IsErrWrongClient is a wrapper.
func IsErrWrongClient(err error) bool {
	return errors.Is(err, mongo.ErrWrongClient)
}

// BulkWriteError is a wrapper.
type BulkWriteError = mongo.BulkWriteError

// BulkWriteException is a wrapper.
type BulkWriteException = mongo.BulkWriteException

// CommandError is a wrapper.
type CommandError = mongo.CommandError

// EncryptionKeyVaultError is a wrapper.
type EncryptionKeyVaultError = mongo.EncryptionKeyVaultError

// ErrMapForOrderedArgument is a wrapper.
type ErrMapForOrderedArgument = mongo.ErrMapForOrderedArgument

// MarshalError is a wrapper.
type MarshalError = mongo.MarshalError

// MongocryptError is a wrapper.
type MongocryptError = mongo.MongocryptError

// MongocryptdError is a wrapper.
type MongocryptdError = mongo.MongocryptdError

// ServerError is a wrapper.
type ServerError = mongo.ServerError

// WriteError is a wrapper.
type WriteError = mongo.WriteError

// WriteErrors is a wrapper.
type WriteErrors = mongo.WriteErrors

// WriteConcernError is a wrapper.
type WriteConcernError = mongo.WriteConcernError

// WriteException is a wrapper.
type WriteException = mongo.WriteException

func wrapError(err error, op string) error {
	err = wrapErrorTemporary(err)
	err = errors.Wrapf(err, "MongoDB %s", op)
	return err
}

func wrapErrorReturn(op string, perr *error, werr func(error) error) {
	err := *perr
	if err != nil {
		err = wrapErrorOptional(err, werr)
		err = wrapError(err, op)
		*perr = err
	}
}

func wrapErrorOptional(err error, werr func(error) error) error {
	if werr != nil {
		err = werr(err)
	}
	return err
}

func wrapErrorValue(err error, key string, val interface{}) error {
	return errors.WithValue(err, "mongo."+key, val)
}

func wrapErrorTemporary(err error) error {
	if !isMongoErrorTemporary(err) {
		err = errors.WithTemporary(err, false)
	}
	return err
}

func isMongoErrorTemporary(err error) bool {
	for _, f := range mongoErrorTemporaryFuncs {
		temporary, ok := f(err)
		if ok {
			return temporary
		}
	}
	return true
}

var mongoErrorTemporaryFuncs = []func(err error) (temporary bool, ok bool){
	func(err error) (temporary bool, ok bool) {
		if IsErrNoDocuments(err) {
			return false, true
		}
		return false, false
	},
	func(err error) (temporary bool, ok bool) {
		var werr MarshalError
		if errors.As(err, &werr) {
			return false, true
		}
		return false, false
	},
	func(err error) (temporary bool, ok bool) {
		var werr *bsoncodec.DecodeError
		if errors.As(err, &werr) {
			return false, true
		}
		return false, false
	},
	func(err error) (temporary bool, ok bool) {
		if IsDuplicateKeyError(err) {
			return false, true
		}
		return false, false
	},
	func(err error) (temporary bool, ok bool) {
		if IsNetworkError(err) {
			return true, true
		}
		return false, false
	},
	func(err error) (temporary bool, ok bool) {
		if IsTimeout(err) {
			return true, true
		}
		return false, false
	},
	func(err error) (temporary bool, ok bool) {
		var werr WriteException
		if errors.As(err, &werr) {
			return !isWriteExceptionCodes(werr, writeErrorNotTemporaryCodes...), true
		}
		return false, false
	},
	func(err error) (temporary bool, ok bool) {
		var werr BulkWriteException
		if errors.As(err, &werr) {
			return !isBulkWriteExceptionCodes(werr, writeErrorNotTemporaryCodes...), true
		}
		return false, false
	},
}

// See https://github.com/mongodb/mongo/blob/master/src/mongo/base/error_codes.yml
var writeErrorNotTemporaryCodes = []int{
	2,     // Bad value
	28,    // Path not viable
	10334, // BSON object too large
	11000, // Duplicate key
	11001, // Duplicate key (legacy)
	12582, // Duplicate key (legacy)
	16837, // Cannot apply $addToSet to non-array field (for older MongoDB version, otherwise it's 2)
	17419, // Resulting document after update is larger than 16777216
	17280, // OBSOLETE KeyTooLong
}

// IsWriteErrorCodes checks if the given error is a write error and the code matches.
//
// This function is not part of the official driver.
// It's used by the current packages and other packages, so it's simpler to define it here.
func IsWriteErrorCodes(err error, codes ...int) bool {
	if werr := (mongo.WriteException{}); errors.As(err, &werr) {
		return isWriteExceptionCodes(werr, codes...)
	}
	if werr := (mongo.BulkWriteException{}); errors.As(err, &werr) {
		return isBulkWriteExceptionCodes(werr, codes...)
	}
	return false
}

func isWriteExceptionCodes(err WriteException, codes ...int) bool {
	for _, e := range err.WriteErrors {
		if isWriteErrorCodes(e, codes...) {
			return true
		}
	}
	return false
}

func isBulkWriteExceptionCodes(err BulkWriteException, codes ...int) bool {
	for _, e := range err.WriteErrors {
		if isWriteErrorCodes(e.WriteError, codes...) {
			return true
		}
	}
	return false
}

func isWriteErrorCodes(err mongo.WriteError, codes ...int) bool {
	for _, c := range codes {
		if err.Code == c {
			return true
		}
	}
	return false
}
