package sqltracing

import (
	"context"
	sqldriver "database/sql/driver"

	"github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/errors"
)

type stmt struct {
	sqldriver.Stmt
	conn  *conn
	query string
}

var _ sqldriver.StmtExecContext = &stmt{}

func (st *stmt) ExecContext(ctx context.Context, args []sqldriver.NamedValue) (r sqldriver.Result, err error) {
	span, spanFinish := st.startSpan(&ctx, "exec", &err)
	defer spanFinish()
	stec, ok := st.Stmt.(sqldriver.StmtExecContext)
	if ok {
		r, err = stec.ExecContext(ctx, args)
	} else {
		var dargs []sqldriver.Value
		dargs, err = namedValueToValue(args)
		if err != nil {
			return nil, errors.Wrap(err, "named value to value")
		}
		return st.Stmt.Exec(dargs) //nolint:staticcheck // We need to handle this legacy interface.
	}
	if err != nil {
		return nil, err
	}
	setSpanTagsResult(span, r)
	return r, nil
}

var _ sqldriver.StmtQueryContext = &stmt{}

func (st *stmt) QueryContext(ctx context.Context, args []sqldriver.NamedValue) (r sqldriver.Rows, err error) {
	_, spanFinish := st.startSpan(&ctx, "query", &err)
	defer spanFinish()
	stqc, ok := st.Stmt.(sqldriver.StmtQueryContext)
	if ok {
		return stqc.QueryContext(ctx, args)
	}
	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, errors.Wrap(err, "named value to value")
	}
	return st.Stmt.Query(dargs) //nolint:staticcheck // We need to handle this legacy interface.
}

var _ sqldriver.ColumnConverter = &stmt{} //nolint:staticcheck // We need to handle this legacy interface.

func (st *stmt) ColumnConverter(idx int) sqldriver.ValueConverter {
	stcc, ok := st.Stmt.(sqldriver.ColumnConverter) //nolint:staticcheck // We need to handle this legacy interface.
	if !ok {
		return sqldriver.DefaultParameterConverter
	}
	return stcc.ColumnConverter(idx)
}

var _ sqldriver.NamedValueChecker = &stmt{}

func (st *stmt) CheckNamedValue(nv *sqldriver.NamedValue) error {
	stcnv, ok := st.Stmt.(sqldriver.NamedValueChecker)
	if ok {
		return stcnv.CheckNamedValue(nv)
	}
	return st.conn.CheckNamedValue(nv)
}

func (st *stmt) startSpan(pctx *context.Context, operationName string, perr *error) (opentracing.Span, closeutils.F) {
	span, spanFinish := startSpan(pctx, "stmt."+operationName, perr)
	st.setSpanTags(span)
	return span, spanFinish
}

func (st *stmt) setSpanTags(span opentracing.Span) {
	st.conn.setSpanTags(span)
	setSpanTagQuery(span, st.query)
}
