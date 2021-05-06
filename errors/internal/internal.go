// Package internal provides internal code for errors.
package internal

type base struct {
	s string
}

// NewBase returns a new base error, just a string.
func NewBase(s string) error {
	return &base{
		s: s,
	}
}

func (err *base) Error() string {
	return err.s
}
