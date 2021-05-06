package mongo

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IndexView is a wrapper.
type IndexView struct {
	indexView
	collection *Collection
}

type indexView = mongo.IndexView

// CreateMany is a wrapper.
func (iv IndexView) CreateMany(ctx context.Context, models []IndexModel, opts ...*options.CreateIndexesOptions) (_ []string, err error) {
	span, spanFinish := iv.startTraceSpan(&ctx, "create_many", &err, func(err error) error {
		err = wrapErrorValue(err, "models.count", len(models))
		return err
	})
	defer spanFinish()
	setTraceSpanTag(span, "models.count", len(models))
	return iv.indexView.CreateMany(ctx, models, opts...)
}

// CreateOne is a wrapper.
func (iv IndexView) CreateOne(ctx context.Context, model IndexModel, opts ...*options.CreateIndexesOptions) (_ string, err error) {
	_, spanFinish := iv.startTraceSpan(&ctx, "create_one", &err, nil)
	defer spanFinish()
	return iv.indexView.CreateOne(ctx, model, opts...)
}

// DropAll is a wrapper.
func (iv IndexView) DropAll(ctx context.Context, opts ...*options.DropIndexesOptions) (_ bson.Raw, err error) {
	_, spanFinish := iv.startTraceSpan(&ctx, "drop_all", &err, nil)
	defer spanFinish()
	return iv.indexView.DropAll(ctx, opts...)
}

// DropOne is a wrapper.
func (iv IndexView) DropOne(ctx context.Context, name string, opts ...*options.DropIndexesOptions) (_ bson.Raw, err error) {
	span, spanFinish := iv.startTraceSpan(&ctx, "drop_one", &err, func(err error) error {
		err = wrapErrorValue(err, "index.name", name)
		return err
	})
	defer spanFinish()
	setTraceSpanTag(span, "index.name", name)
	return iv.indexView.DropOne(ctx, name, opts...)
}

// List is a wrapper.
func (iv IndexView) List(ctx context.Context, opts ...*options.ListIndexesOptions) (_ *Cursor, err error) {
	_, spanFinish := iv.startTraceSpan(&ctx, "list", &err, nil)
	defer spanFinish()
	cu, err := iv.indexView.List(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &Cursor{
		cursor:     cu,
		collection: iv.collection,
	}, nil
}

// ListSpecifications is a wrapper.
func (iv IndexView) ListSpecifications(ctx context.Context, opts ...*options.ListIndexesOptions) (_ []*IndexSpecification, err error) {
	_, spanFinish := iv.startTraceSpan(&ctx, "list_specifications", &err, nil)
	defer spanFinish()
	return iv.indexView.ListSpecifications(ctx, opts...)
}

func (iv IndexView) startTraceSpan(pctx *context.Context, op string, perr *error, werr func(error) error) (opentracing.Span, closeutils.F) {
	span, spanFinish := startTraceSpan(pctx, "index_view."+op, perr, func(err error) error {
		err = iv.wrapErrorValues(err)
		err = wrapErrorOptional(err, werr)
		return err
	})
	iv.setTraceSpanTags(span)
	return span, spanFinish
}

func (iv IndexView) setTraceSpanTags(span opentracing.Span) {
	iv.collection.setTraceSpanTags(span)
}

func (iv IndexView) wrapErrorValues(err error) error {
	err = iv.collection.wrapErrorValues(err)
	return err
}

// IndexModel is a wrapper.
type IndexModel = mongo.IndexModel

// IndexSpecification is a wrapper.
type IndexSpecification = mongo.IndexSpecification
