package mongo

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database is a wrapper.
type Database struct {
	*database
	client *Client
}

type database = mongo.Database

// Aggregate is a wrapper.
func (db *Database) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (_ *Cursor, err error) {
	span, spanFinish := db.startTraceSpan(&ctx, "aggregate", &err, func(err error) error {
		if pipeline != nil {
			err = wrapErrorValue(err, "pipeline", pipeline)
		}
		return err
	})
	defer spanFinish()
	if pipeline != nil {
		setTraceSpanTagJSON(span, "pipeline", pipeline)
	}
	cu, err := db.database.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}
	return &Cursor{
		cursor:   cu,
		database: db,
		pipeline: pipeline,
	}, nil
}

// Client is a wrapper.
func (db *Database) Client() *Client {
	return db.client
}

// Collection is a wrapper.
func (db *Database) Collection(name string, opts ...*options.CollectionOptions) *Collection {
	return &Collection{
		collection: db.database.Collection(name, opts...),
		database:   db,
	}
}

// CreateCollection is a wrapper.
func (db *Database) CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) (err error) {
	span, spanFinish := db.startTraceSpan(&ctx, "create_collection", &err, func(err error) error {
		err = wrapErrorValue(err, "collection.name", name)
		return err
	})
	defer spanFinish()
	setTraceSpanTag(span, "collection.name", name)
	return db.database.CreateCollection(ctx, name, opts...)
}

// CreateView is a wrapper.
func (db *Database) CreateView(ctx context.Context, viewName, viewOn string, pipeline interface{}, opts ...*options.CreateViewOptions) (err error) {
	span, spanFinish := db.startTraceSpan(&ctx, "create_view", &err, func(err error) error {
		err = wrapErrorValue(err, "view.name", viewName)
		err = wrapErrorValue(err, "view.on", viewOn)
		return err
	})
	defer spanFinish()
	setTraceSpanTag(span, "view.name", viewName)
	setTraceSpanTag(span, "view.on", viewOn)
	return db.database.CreateView(ctx, viewName, viewOn, pipeline, opts...)
}

// Drop is a wrapper.
func (db *Database) Drop(ctx context.Context) (err error) {
	_, spanFinish := db.startTraceSpan(&ctx, "drop", &err, nil)
	defer spanFinish()
	return db.database.Drop(ctx)
}

// ListCollectionNames is a wrapper.
func (db *Database) ListCollectionNames(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) (_ []string, err error) {
	span, spanFinish := db.startTraceSpan(&ctx, "list_collection_names", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	return db.database.ListCollectionNames(ctx, filter, opts...)
}

// ListCollectionSpecifications is a wrapper.
func (db *Database) ListCollectionSpecifications(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) (_ []*CollectionSpecification, err error) {
	span, spanFinish := db.startTraceSpan(&ctx, "list_collection_specifications", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	return db.database.ListCollectionSpecifications(ctx, filter, opts...)
}

// ListCollections is a wrapper.
func (db *Database) ListCollections(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) (_ *Cursor, err error) {
	span, spanFinish := db.startTraceSpan(&ctx, "list_collections", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	cu, err := db.database.ListCollections(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return &Cursor{
		cursor:   cu,
		database: db,
		filter:   filter,
	}, nil
}

// RunCommand is a wrapper.
func (db *Database) RunCommand(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *SingleResult {
	var err error
	span, spanFinish := db.startTraceSpan(&ctx, "run_command", &err, func(err error) error {
		if runCommand != nil {
			err = wrapErrorValue(err, "run_command", runCommand)
		}
		return err
	})
	defer spanFinish()
	if runCommand != nil {
		setTraceSpanTagJSON(span, "run_command", runCommand)
	}
	res := db.database.RunCommand(ctx, runCommand, opts...)
	err = res.Err()
	return &SingleResult{
		singleResult: res,
		database:     db,
		runCommand:   runCommand,
	}
}

// RunCommandCursor is a wrapper.
func (db *Database) RunCommandCursor(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) (_ *Cursor, err error) {
	span, spanFinish := db.startTraceSpan(&ctx, "run_command_cursor", &err, func(err error) error {
		if runCommand != nil {
			err = wrapErrorValue(err, "run_command", runCommand)
		}
		return err
	})
	defer spanFinish()
	if runCommand != nil {
		setTraceSpanTagJSON(span, "run_command", runCommand)
	}
	cu, err := db.database.RunCommandCursor(ctx, runCommand, opts...)
	if err != nil {
		return nil, err
	}
	return &Cursor{
		cursor:     cu,
		database:   db,
		runCommand: runCommand,
	}, nil
}

// Watch is a wrapper.
func (db *Database) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (_ *ChangeStream, err error) {
	span, spanFinish := db.startTraceSpan(&ctx, "watch", &err, func(err error) error {
		if pipeline != nil {
			err = wrapErrorValue(err, "pipeline", pipeline)
		}
		return err
	})
	defer spanFinish()
	if pipeline != nil {
		setTraceSpanTagJSON(span, "pipeline", pipeline)
	}
	cs, err := db.database.Watch(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}
	return &ChangeStream{
		changeStream: cs,
		database:     db,
		pipeline:     pipeline,
	}, nil
}

func (db *Database) startTraceSpan(pctx *context.Context, op string, perr *error, werr func(error) error) (opentracing.Span, closeutils.F) {
	span, spanFinish := startTraceSpan(pctx, "database."+op, perr, func(err error) error {
		err = db.wrapErrorValues(err)
		err = wrapErrorOptional(err, werr)
		return err
	})
	db.setTraceSpanTags(span)
	return span, spanFinish
}

func (db *Database) setTraceSpanTags(span opentracing.Span) {
	db.client.setTraceSpanTags(span)
	setTraceSpanTag(span, "database.name", db.Name())
}

func (db *Database) wrapErrorValues(err error) error {
	err = db.client.wrapErrorValues(err)
	err = wrapErrorValue(err, "database.name", db.Name())
	return err
}
