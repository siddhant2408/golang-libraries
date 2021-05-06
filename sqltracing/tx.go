package sqltracing

import (
	sqldriver "database/sql/driver"

	"github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
)

type tx struct {
	sqldriver.Tx
	span       opentracing.Span
	spanFinish closeutils.F
}

func (t *tx) Commit() error {
	t.setSpanTagResult("commit")
	t.spanFinish()
	return t.Tx.Commit()
}

func (t *tx) Rollback() error {
	t.setSpanTagResult("rollback")
	t.spanFinish()
	return t.Tx.Rollback()
}

func (t *tx) setSpanTag(key string, value interface{}) {
	setSpanTag(t.span, "transaction."+key, value)
}

func (t *tx) setSpanTagResult(value string) {
	t.setSpanTag("result", value)
}
