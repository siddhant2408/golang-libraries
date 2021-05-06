package sibvalid

import (
	"regexp"

	"github.com/siddhant2408/golang-libraries/errors"
)

var phone = regexp.MustCompile(`^[1-9]\d{5,15}$`)

// Phone validates a phone number.
func Phone(s string) error {
	ok := phone.MatchString(s)
	if !ok {
		return errors.New("regexp")
	}
	return nil
}
