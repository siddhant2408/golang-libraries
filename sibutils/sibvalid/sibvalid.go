// Package sibvalid provides validation according to rules.
package sibvalid

import (
	"unicode/utf8"

	"github.com/siddhant2408/golang-libraries/errors"
)

func checkRunesCountMax(s string, max int) error {
	l := utf8.RuneCountInString(s)
	if l > max {
		return errors.Newf("greater than the maximum allowed: got %d, max %d", l, max)
	}
	return nil
}
