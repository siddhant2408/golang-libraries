package sqltracing

import (
	"context"
	sqldriver "database/sql/driver"

	"github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/errors"
)

type conn struct {
	sqldriver.Conn
	connector *connector
}

var _ sqldriver.ConnPrepareContext = &conn{}

func (c *conn) PrepareContext(ctx context.Context, query string) (st sqldriver.Stmt, err error) {
	span, spanFinish := c.startSpan(&ctx, "prepare", &err)
	defer spanFinish()
	setSpanTagQuery(span, query)
	cpc, ok := c.Conn.(sqldriver.ConnPrepareContext)
	if ok {
		st, err = cpc.PrepareContext(ctx, query)
	} else {
		st, err = c.Conn.Prepare(query)
	}
	if err != nil {
		return nil, err
	}
	return &stmt{
		Stmt:  st,
		conn:  c,
		query: query,
	}, nil
}

var _ sqldriver.ConnBeginTx = &conn{}

func (c *conn) BeginTx(ctx context.Context, opts sqldriver.TxOptions) (t sqldriver.Tx, err error) {
	span, spanFinish := c.startSpan(&ctx, "transaction", &err)
	cbt, ok := c.Conn.(sqldriver.ConnBeginTx)
	if ok {
		t, err = cbt.BeginTx(ctx, opts)
	} else {
		t, err = c.Conn.Begin() //nolint:staticcheck // We need to handle this legacy interface.
	}
	if err != nil {
		return nil, err
	}
	return &tx{
		Tx:         t,
		span:       span,
		spanFinish: spanFinish,
	}, nil
}

var _ sqldriver.ExecerContext = &conn{}

func (c *conn) ExecContext(ctx context.Context, query string, args []sqldriver.NamedValue) (r sqldriver.Result, err error) {
	cec, _ := c.Conn.(sqldriver.ExecerContext)
	ce, _ := c.Conn.(sqldriver.Execer) //nolint:staticcheck // We need to handle this legacy interface.
	if cec == nil && ce == nil {
		return nil, sqldriver.ErrSkip
	}
	span, spanFinish := c.startSpan(&ctx, "exec", &err)
	defer spanFinish()
	setSpanTagQuery(span, query)
	if cec != nil {
		r, err = cec.ExecContext(ctx, query, args)
	} else {
		var dargs []sqldriver.Value
		dargs, err = namedValueToValue(args)
		if err != nil {
			return nil, errors.Wrap(err, "named value to value")
		}
		r, err = ce.Exec(query, dargs)
	}
	if err != nil {
		return nil, err
	}
	setSpanTagsResult(span, r)
	return r, nil
}

var _ sqldriver.QueryerContext = &conn{}

func (c *conn) QueryContext(ctx context.Context, query string, args []sqldriver.NamedValue) (r sqldriver.Rows, err error) {
	cqc, _ := c.Conn.(sqldriver.QueryerContext)
	cq, _ := c.Conn.(sqldriver.Queryer) //nolint:staticcheck // We need to handle this legacy interface.
	if cqc == nil && cq == nil {
		return nil, sqldriver.ErrSkip
	}
	span, spanFinish := c.startSpan(&ctx, "query", &err)
	defer spanFinish()
	setSpanTagQuery(span, query)
	if cqc != nil {
		return cqc.QueryContext(ctx, query, args)
	}
	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, errors.Wrap(err, "named value to value")
	}
	return cq.Query(query, dargs)
}

var _ sqldriver.NamedValueChecker = &conn{}

func (c *conn) CheckNamedValue(nv *sqldriver.NamedValue) error {
	ccnv, ok := c.Conn.(sqldriver.NamedValueChecker)
	if !ok {
		return sqldriver.ErrSkip
	}
	return ccnv.CheckNamedValue(nv)
}

var _ sqldriver.Pinger = &conn{}

func (c *conn) Ping(ctx context.Context) (err error) {
	p, ok := c.Conn.(sqldriver.Pinger)
	if !ok {
		return nil
	}
	_, spanFinish := c.startSpan(&ctx, "ping", &err)
	defer spanFinish()
	return p.Ping(ctx)
}

var _ sqldriver.SessionResetter = &conn{}

func (c *conn) ResetSession(ctx context.Context) (err error) {
	crs, ok := c.Conn.(sqldriver.SessionResetter)
	if !ok {
		return nil
	}
	_, spanFinish := c.startSpan(&ctx, "reset_session", &err)
	defer spanFinish()
	return crs.ResetSession(ctx)
}

var _ sqldriver.Validator = &conn{}

func (c *conn) IsValid() bool {
	cis, ok := c.Conn.(sqldriver.Validator)
	if !ok {
		return true
	}
	return cis.IsValid()
}

func (c *conn) startSpan(pctx *context.Context, operationName string, perr *error) (opentracing.Span, closeutils.F) {
	span, spanFinish := startSpan(pctx, "conn."+operationName, perr)
	c.setSpanTags(span)
	return span, spanFinish
}

func (c *conn) setSpanTags(span opentracing.Span) {
	c.connector.setSpanTags(span)
}
