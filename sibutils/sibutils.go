// Package sibutils contains utils specific to sib.
package sibutils

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/siddhant2408/golang-libraries/errors"
)

// Encode encodes sib reference string.
func Encode(v []int64) (string, error) {
	if len(v) == 0 {
		return "", nil
	}
	vs, err := intSliceToString(v)
	if err != nil {
		return "", errors.Wrap(err, "convert ints to strings")
	}
	s := strings.Join(vs, "a")
	s, err = convertBase(s, 12, 36)
	if err != nil {
		return "", errors.Wrap(err, "convert base 12 to 36")
	}
	return s, nil
}

func intSliceToString(vi []int64) ([]string, error) {
	vs := make([]string, len(vi))
	for i, n := range vi {
		if n < 0 {
			return nil, errors.Newf("negative value \"%d\" is not supported", n)
		}
		vs[i] = strconv.FormatInt(n, 10)
	}
	return vs, nil
}

// Decode decodes sib reference string.
func Decode(ref string) ([]int64, error) {
	if ref == "" {
		return nil, nil
	}
	s, err := convertBase(ref, 36, 12)
	if err != nil {
		return nil, errors.Wrap(err, "convert base 36 to 12")
	}
	vs := strings.Split(s, "a")
	vi, err := stringSliceToInt(vs)
	if err != nil {
		return nil, errors.Wrap(err, "convert strings to ints")
	}
	return vi, nil
}

func stringSliceToInt(vs []string) ([]int64, error) {
	vi := make([]int64, len(vs))
	for i, s := range vs {
		if s != "" {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
			vi[i] = n
		}
	}
	return vi, nil
}

func convertBase(s string, from, to int) (string, error) {
	i := new(big.Int)
	_, ok := i.SetString(s, from)
	if !ok {
		return "", errors.Newf("%q is not a valid base %d integer", s, from)
	}
	return i.Text(to), nil
}
