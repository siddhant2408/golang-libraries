package sqltracing

import (
	"context"
	"database/sql"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

// RowScan traces the execution of Row.Scan.
func RowScan(ctx context.Context, row *sql.Row, dest ...interface{}) (err error) {
	span, spanFinish := startSpan(&ctx, "row.scan", nil)
	defer spanFinish()
	err = row.Scan(dest...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			setSpanTag(span, "row.found", false)
		} else {
			tracingutils.SetSpanError(span, err)
		}
		return errors.Wrap(err, "")
	}
	setSpanTag(span, "row.found", true)
	return nil
}

// RowsIterate traces the iteration of Rows.
//
// It closes the rows and check the error.
func RowsIterate(ctx context.Context, rows *sql.Rows, f func(context.Context, *sql.Rows) error) (err error) {
	span, spanFinish := startSpan(&ctx, "rows.iterate", &err)
	defer spanFinish()
	defer rows.Close() //nolint:errcheck
	count := 0
	defer func() {
		setSpanTag(span, "rows.count", count)
	}()
	for rows.Next() {
		count++
		err = f(ctx, rows)
		if err != nil {
			return errors.Wrap(err, "row")
		}
	}
	err = rows.Err()
	if err != nil {
		return errors.Wrap(err, "rows")
	}
	return nil
}
