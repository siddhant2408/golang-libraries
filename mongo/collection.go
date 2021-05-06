package mongo

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection is a wrapper.
type Collection struct {
	*collection
	database *Database
}

type collection = mongo.Collection

// Aggregate is a wrapper.
func (coll *Collection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (_ *Cursor, err error) {
	span, spanFinish := coll.startTraceSpan(&ctx, "aggregate", &err, func(err error) error {
		if pipeline != nil {
			err = wrapErrorValue(err, "pipeline", pipeline)
		}
		return err
	})
	defer spanFinish()
	if pipeline != nil {
		setTraceSpanTagJSON(span, "pipeline", pipeline)
	}
	cu, err := coll.collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}
	return &Cursor{
		cursor:     cu,
		collection: coll,
		pipeline:   pipeline,
	}, nil
}

// BulkWrite is a wrapper.
func (coll *Collection) BulkWrite(ctx context.Context, models []WriteModel, opts ...*options.BulkWriteOptions) (_ *BulkWriteResult, err error) {
	span, spanFinish := coll.startTraceSpan(&ctx, "bulk_write", &err, func(err error) error {
		err = wrapErrorValue(err, "models.count", len(models))
		return err
	})
	defer spanFinish()
	setTraceSpanTag(span, "models.count", len(models))
	return coll.collection.BulkWrite(ctx, models, opts...)
}

// Clone is a wrapper.
func (coll *Collection) Clone(opts ...*options.CollectionOptions) (_ *Collection, err error) {
	defer coll.wrapErrorReturn("clone", &err, nil)
	c, err := coll.collection.Clone(opts...)
	if err != nil {
		return nil, err
	}
	return &Collection{
		collection: c,
		database:   coll.database,
	}, nil
}

// CountDocuments is a wrapper.
func (coll *Collection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (_ int64, err error) {
	span, spanFinish := coll.startTraceSpan(&ctx, "count_documents", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	return coll.collection.CountDocuments(ctx, filter, opts...)
}

// Database is a wrapper.
func (coll *Collection) Database() *Database {
	return coll.database
}

// DeleteMany is a wrapper.
func (coll *Collection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (_ *DeleteResult, err error) {
	span, spanFinish := coll.startTraceSpan(&ctx, "delete_many", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	return coll.collection.DeleteMany(ctx, filter, opts...)
}

// DeleteOne is a wrapper.
func (coll *Collection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (_ *DeleteResult, err error) {
	span, spanFinish := coll.startTraceSpan(&ctx, "delete_one", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	return coll.collection.DeleteOne(ctx, filter, opts...)
}

// Distinct is a wrapper.
func (coll *Collection) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) (_ []interface{}, err error) {
	span, spanFinish := coll.startTraceSpan(&ctx, "distinct", &err, func(err error) error {
		err = wrapErrorValue(err, "field.name", fieldName)
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	setTraceSpanTag(span, "field.name", fieldName)
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	return coll.collection.Distinct(ctx, fieldName, filter, opts...)
}

// Drop is a wrapper.
func (coll *Collection) Drop(ctx context.Context) (err error) {
	_, spanFinish := coll.startTraceSpan(&ctx, "drop", &err, nil)
	defer spanFinish()
	return coll.collection.Drop(ctx)
}

// EstimatedDocumentCount is a wrapper.
func (coll *Collection) EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (_ int64, err error) {
	_, spanFinish := coll.startTraceSpan(&ctx, "estimated_document_count", &err, nil)
	defer spanFinish()
	return coll.collection.EstimatedDocumentCount(ctx, opts...)
}

// Find is a wrapper.
func (coll *Collection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (_ *Cursor, err error) {
	span, spanFinish := coll.startTraceSpan(&ctx, "find", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	cu, err := coll.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return &Cursor{
		cursor:     cu,
		collection: coll,
		filter:     filter,
	}, nil
}

// FindOne is a wrapper.
func (coll *Collection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *SingleResult {
	var err error
	span, spanFinish := coll.startTraceSpan(&ctx, "find_one", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	res := coll.collection.FindOne(ctx, filter, opts...)
	found := true
	err = res.Err()
	if err != nil && IsErrNoDocuments(err) {
		err = nil
		found = false
	}
	if err == nil {
		setTraceSpanTag(span, "found", found)
	}
	return &SingleResult{
		singleResult: res,
		collection:   coll,
		filter:       filter,
	}
}

// FindOneAndDelete is a wrapper.
func (coll *Collection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *SingleResult {
	var err error
	span, spanFinish := coll.startTraceSpan(&ctx, "find_one_and_delete", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	res := coll.collection.FindOneAndDelete(ctx, filter, opts...)
	found := true
	err = res.Err()
	if err != nil && IsErrNoDocuments(err) {
		err = nil
		found = false
	}
	if err == nil {
		setTraceSpanTag(span, "found", found)
	}
	return &SingleResult{
		singleResult: res,
		collection:   coll,
		filter:       filter,
	}
}

// FindOneAndReplace is a wrapper.
func (coll *Collection) FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *SingleResult {
	var err error
	span, spanFinish := coll.startTraceSpan(&ctx, "find_one_and_replace", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	res := coll.collection.FindOneAndReplace(ctx, filter, replacement, opts...)
	found := true
	err = res.Err()
	if err != nil && IsErrNoDocuments(err) {
		err = nil
		found = false
	}
	if err == nil {
		setTraceSpanTag(span, "found", found)
	}
	return &SingleResult{
		singleResult: res,
		collection:   coll,
		filter:       filter,
	}
}

// FindOneAndUpdate is a wrapper.
func (coll *Collection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *SingleResult {
	var err error
	span, spanFinish := coll.startTraceSpan(&ctx, "find_one_and_update", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	res := coll.collection.FindOneAndUpdate(ctx, filter, update, opts...)
	found := true
	err = res.Err()
	if err != nil && IsErrNoDocuments(err) {
		err = nil
		found = false
	}
	if err == nil {
		setTraceSpanTag(span, "found", found)
	}
	return &SingleResult{
		singleResult: res,
		collection:   coll,
		filter:       filter,
	}
}

// Indexes is a wrapper.
func (coll *Collection) Indexes() IndexView {
	return IndexView{
		indexView:  coll.collection.Indexes(),
		collection: coll,
	}
}

// InsertMany is a wrapper.
func (coll *Collection) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (_ *InsertManyResult, err error) {
	span, spanFinish := coll.startTraceSpan(&ctx, "insert_many", &err, func(err error) error {
		err = wrapErrorValue(err, "documents.count", len(documents))
		return err
	})
	defer spanFinish()
	setTraceSpanTag(span, "documents.count", len(documents))
	return coll.collection.InsertMany(ctx, documents, opts...)
}

// InsertOne is a wrapper.
func (coll *Collection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (_ *InsertOneResult, err error) {
	_, spanFinish := coll.startTraceSpan(&ctx, "insert_one", &err, nil)
	defer spanFinish()
	return coll.collection.InsertOne(ctx, document, opts...)
}

// ReplaceOne is a wrapper.
func (coll *Collection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (_ *UpdateResult, err error) {
	span, spanFinish := coll.startTraceSpan(&ctx, "replace_one", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	return coll.collection.ReplaceOne(ctx, filter, replacement, opts...)
}

// UpdateByID is a wrapper.
func (coll *Collection) UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (_ *UpdateResult, err error) {
	span, spanFinish := coll.startTraceSpan(&ctx, "update_by_id", &err, func(err error) error {
		if id != nil {
			err = wrapErrorValue(err, "id", id)
		}
		return err
	})
	defer spanFinish()
	if id != nil {
		setTraceSpanTagJSON(span, "id", fmt.Sprint(id))
	}
	return coll.collection.UpdateByID(ctx, id, update, opts...)
}

// UpdateMany is a wrapper.
func (coll *Collection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (_ *UpdateResult, err error) {
	span, spanFinish := coll.startTraceSpan(&ctx, "update_many", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	return coll.collection.UpdateMany(ctx, filter, update, opts...)
}

// UpdateOne is a wrapper.
func (coll *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (_ *UpdateResult, err error) {
	span, spanFinish := coll.startTraceSpan(&ctx, "update_one", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	return coll.collection.UpdateOne(ctx, filter, update, opts...)
}

// Watch is a wrapper.
func (coll *Collection) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (_ *ChangeStream, err error) {
	span, spanFinish := coll.startTraceSpan(&ctx, "watch", &err, func(err error) error {
		if pipeline != nil {
			err = wrapErrorValue(err, "pipeline", pipeline)
		}
		return err
	})
	defer spanFinish()
	if pipeline != nil {
		setTraceSpanTagJSON(span, "pipeline", pipeline)
	}
	cs, err := coll.collection.Watch(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}
	return &ChangeStream{
		changeStream: cs,
		collection:   coll,
		pipeline:     pipeline,
	}, nil
}

func (coll *Collection) startTraceSpan(pctx *context.Context, op string, perr *error, werr func(error) error) (opentracing.Span, closeutils.F) {
	span, spanFinish := startTraceSpan(pctx, "collection."+op, perr, func(err error) error {
		err = coll.wrapErrorValues(err)
		err = wrapErrorOptional(err, werr)
		return err
	})
	coll.setTraceSpanTags(span)
	return span, spanFinish
}

func (coll *Collection) setTraceSpanTags(span opentracing.Span) {
	coll.database.setTraceSpanTags(span)
	setTraceSpanTag(span, "collection.name", coll.Name())
}

func (coll *Collection) wrapErrorReturn(op string, perr *error, werr func(error) error) {
	wrapErrorReturn("collection."+op, perr, func(err error) error {
		err = coll.wrapErrorValues(err)
		err = wrapErrorOptional(err, werr)
		return err
	})
}

func (coll *Collection) wrapErrorValues(err error) error {
	err = coll.database.wrapErrorValues(err)
	err = wrapErrorValue(err, "collection.name", coll.Name())
	return err
}

// CollectionSpecification is a wrapper.
type CollectionSpecification = mongo.CollectionSpecification
