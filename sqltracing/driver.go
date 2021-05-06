package sqltracing

import (
	sqldriver "database/sql/driver"
)

type driver struct {
	sqldriver.Driver
}

// WrapDriver adds tracing to a Driver.
func WrapDriver(drv sqldriver.Driver) sqldriver.Driver {
	return &driver{
		Driver: drv,
	}
}

var _ sqldriver.DriverContext = &driver{}

func (d *driver) OpenConnector(dsn string) (cr sqldriver.Connector, err error) {
	dc, ok := d.Driver.(sqldriver.DriverContext)
	if ok {
		cr, err = dc.OpenConnector(dsn)
		if err != nil {
			return nil, err
		}
	} else {
		cr = &dsnConnector{
			driver: d.Driver,
			dsn:    dsn,
		}
	}
	return &connector{
		Connector: cr,
		driver:    d,
		dsn:       dsn,
	}, nil
}
