package sqltracing

import (
	"context"
	sqldriver "database/sql/driver"

	"github.com/opentracing/opentracing-go"
	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

const (
	// ExternalServiceName is the external service name.
	ExternalServiceName = "go-sql"
)

func startSpan(pctx *context.Context, operationName string, perr *error) (opentracing.Span, closeutils.F) {
	span, cl := tracingutils.StartChildSpan(pctx, "sql."+operationName, perr)
	tracingutils.SetSpanServiceName(span, ExternalServiceName)
	tracingutils.SetSpanType(span, tracingutils.SpanTypeSQL)
	opentracing_ext.SpanKindRPCClient.Set(span)
	return span, cl
}

func setSpanTag(span opentracing.Span, key string, value interface{}) {
	span.SetTag("sql."+key, value)
}

func setSpanTagQuery(span opentracing.Span, query string) {
	setSpanTag(span, "query", query)
}

func setSpanTagsResult(span opentracing.Span, r sqldriver.Result) {
	lastInsertID, err := r.LastInsertId()
	if err == nil {
		setSpanTag(span, "result.last_insert_id", lastInsertID)
	}
	rowsAffected, err := r.RowsAffected()
	if err == nil {
		setSpanTag(span, "result.rows_affected", rowsAffected)
	}
}
