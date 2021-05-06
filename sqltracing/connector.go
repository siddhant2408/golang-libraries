package sqltracing

import (
	"context"
	sqldriver "database/sql/driver"

	"github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
)

type connector struct {
	sqldriver.Connector
	driver *driver
	dsn    string
}

// WrapConnector adds tracing to a Connector.
func WrapConnector(cr sqldriver.Connector, dsn string) sqldriver.Connector {
	return &connector{
		Connector: cr,
		driver: &driver{
			Driver: cr.Driver(),
		},
		dsn: dsn,
	}
}

func (cr *connector) Connect(ctx context.Context) (c sqldriver.Conn, err error) {
	span, spanFinish := cr.startSpan(&ctx, "connect", &err)
	defer spanFinish()
	cr.setSpanTags(span)
	c, err = cr.Connector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return &conn{
		Conn:      c,
		connector: cr,
	}, nil
}

func (cr *connector) startSpan(pctx *context.Context, operationName string, perr *error) (opentracing.Span, closeutils.F) {
	span, spanFinish := startSpan(pctx, "connector."+operationName, perr)
	cr.setSpanTags(span)
	return span, spanFinish
}

func (cr *connector) setSpanTags(span opentracing.Span) {
	setSpanTag(span, "dsn", cr.dsn)
}

func (cr *connector) Driver() sqldriver.Driver {
	return cr.driver
}

type dsnConnector struct {
	driver sqldriver.Driver
	dsn    string
}

func (cr *dsnConnector) Connect(_ context.Context) (sqldriver.Conn, error) {
	return cr.driver.Open(cr.dsn)
}

func (cr *dsnConnector) Driver() sqldriver.Driver {
	return cr.driver
}
