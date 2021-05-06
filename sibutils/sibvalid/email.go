package sibvalid

import (
	"regexp"
	"unicode/utf8"

	"github.com/siddhant2408/golang-libraries/errors"
)

const (
	emailMaxLen     = 255
	emailRegexpStrX = `[#&*\/=?^{!}~'_\pL0-9-\+]`
	emailRegexpStrY = `[_\pL0-9-]`
	emailRegexpStr  = `(?i)^` + emailRegexpStrX + `+(\.` + emailRegexpStrX + `+)*\.?@(` + emailRegexpStrY + `+(\.` + emailRegexpStrY + `+)*\.)?[\pL0-9-]*[\pL0-9]$`
)

var (
	emailRegexp = regexp.MustCompile(emailRegexpStr)
	emailChecks = []struct {
		msg   string
		check func(string) error
	}{
		{
			msg:   "UTF-8",
			check: checkEmailUTF8Valid,
		},
		{
			msg:   "runes count",
			check: checkEmailRunesCount,
		},
		{
			msg:   "regexp",
			check: checkEmailRegexp,
		},
	}
)

// Email validates an email address.
func Email(s string) error {
	for _, v := range emailChecks {
		err := v.check(s)
		if err != nil {
			return errors.Wrap(err, v.msg)
		}
	}
	return nil
}

func checkEmailUTF8Valid(s string) error {
	ok := utf8.ValidString(s)
	if !ok {
		return errors.New("invalid")
	}
	return nil
}

func checkEmailRunesCount(s string) error {
	return checkRunesCountMax(s, emailMaxLen)
}

func checkEmailRegexp(s string) error {
	ok := emailRegexp.MatchString(s)
	if !ok {
		return errors.New("no match")
	}
	return nil
}
