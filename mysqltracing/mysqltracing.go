// Package mysqltracing provides tracing for MySQL.
//
// Importing this package registers a SQL driver named "mysql-tracing".
package mysqltracing

import (
	"database/sql"
	sqldriver "database/sql/driver"

	"github.com/go-sql-driver/mysql"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/sqltracing"
)

// DriverName is the driver name.
const DriverName = "mysql-tracing"

func init() {
	d := sqltracing.WrapDriver(&mysql.MySQLDriver{})
	sql.Register(DriverName, d)
}

// NewConnector returns a new Connector with tracing.
func NewConnector(cfg *mysql.Config) (sqldriver.Connector, error) {
	cr, err := mysql.NewConnector(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	dsn := cfg.FormatDSN()
	cr = sqltracing.WrapConnector(cr, dsn)
	return cr, nil
}
