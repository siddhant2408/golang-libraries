package sqltracing

import (
	sqldriver "database/sql/driver"

	"github.com/siddhant2408/golang-libraries/errors"
)

func namedValueToValue(named []sqldriver.NamedValue) ([]sqldriver.Value, error) {
	dargs := make([]sqldriver.Value, len(named))
	for n, param := range named {
		if len(param.Name) > 0 {
			return nil, errors.New("sql: driver does not support the use of Named Parameters")
		}
		dargs[n] = param.Value
	}
	return dargs, nil
}
