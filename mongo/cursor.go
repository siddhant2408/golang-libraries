package mongo

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

// Cursor is a wrapper.
type Cursor struct {
	*cursor
	database   *Database
	collection *Collection
	filter     interface{}
	pipeline   interface{}
	runCommand interface{}
}

type cursor = mongo.Cursor

// All is a wrapper.
func (c *Cursor) All(ctx context.Context, results interface{}) (err error) {
	_, spanFinish := c.startTraceSpan(&ctx, "all", &err, nil)
	defer spanFinish()
	return c.cursor.All(ctx, results)
}

// Close is a wrapper.
func (c *Cursor) Close(ctx context.Context) (err error) {
	defer c.wrapErrorReturn("close", &err, nil)
	return c.cursor.Close(ctx)
}

// Decode is a wrapper.
func (c *Cursor) Decode(val interface{}) (err error) {
	defer c.wrapErrorReturn("decode", &err, nil)
	return c.cursor.Decode(val)
}

// Err is a wrapper.
func (c *Cursor) Err() (err error) {
	defer c.wrapErrorReturn("err", &err, nil)
	return c.cursor.Err()
}

func (c *Cursor) startTraceSpan(pctx *context.Context, op string, perr *error, werr func(error) error) (opentracing.Span, closeutils.F) {
	span, spanFinish := startTraceSpan(pctx, "cursor."+op, perr, func(err error) error {
		err = c.wrapErrorValues(err)
		err = wrapErrorOptional(err, werr)
		return err
	})
	c.setTraceSpanTags(span)
	return span, spanFinish
}

func (c *Cursor) setTraceSpanTags(span opentracing.Span) {
	if c.database != nil {
		c.database.setTraceSpanTags(span)
	}
	if c.collection != nil {
		c.collection.setTraceSpanTags(span)
	}
	if c.filter != nil {
		setTraceSpanTagJSON(span, "filter", c.filter)
	}
	if c.pipeline != nil {
		setTraceSpanTagJSON(span, "pipeline", c.pipeline)
	}
	if c.runCommand != nil {
		setTraceSpanTagJSON(span, "run_command", c.runCommand)
	}
}

func (c *Cursor) wrapErrorReturn(op string, perr *error, werr func(error) error) {
	wrapErrorReturn("cursor."+op, perr, func(err error) error {
		err = c.wrapErrorValues(err)
		err = wrapErrorOptional(err, werr)
		return err
	})
}

func (c *Cursor) wrapErrorValues(err error) error {
	if c.database != nil {
		err = c.database.wrapErrorValues(err)
	}
	if c.collection != nil {
		err = c.collection.wrapErrorValues(err)
	}
	if c.filter != nil {
		err = wrapErrorValue(err, "filter", c.filter)
	}
	if c.pipeline != nil {
		err = wrapErrorValue(err, "pipeline", c.pipeline)
	}
	if c.runCommand != nil {
		err = wrapErrorValue(err, "run_command", c.runCommand)
	}
	return err
}

// BatchCursorFromCursor is a wrapper.
func BatchCursorFromCursor(c *mongo.Cursor) *driver.BatchCursor {
	return mongo.BatchCursorFromCursor(c) //nolint:staticcheck // This is a wrapper.
}
