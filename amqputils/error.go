package amqputils

import (
	"github.com/siddhant2408/golang-libraries/errors"
)

func wrapErrorValue(err error, key string, val interface{}) error {
	return errors.WithValue(err, "amqp."+key, val)
}

func wrapErrorValueBody(err error, body []byte) error {
	return wrapErrorValue(err, "body", bodyTruncateConvert(body))
}
