package mongo

import (
	"context"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client is a wrapper.
type Client struct {
	*client
	options *options.ClientOptions
}

type client = mongo.Client

// Connect is a wrapper.
func Connect(ctx context.Context, opts ...*options.ClientOptions) (c *Client, err error) {
	_, spanFinish := startTraceSpan(&ctx, "connect", &err, nil)
	defer spanFinish()
	c, err = NewClient(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "new client")
	}
	err = c.Connect(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	return c, nil
}

// NewClient is a wrapper.
func NewClient(opts ...*options.ClientOptions) (*Client, error) {
	mopts := options.MergeClientOptions(opts...)
	c, err := mongo.NewClient(opts...)
	if err != nil {
		err = wrapClientOptionsErrorValues(mopts, err)
		return nil, err
	}
	return &Client{
		client:  c,
		options: mopts,
	}, nil
}

// Connect is a wrapper.
func (c *Client) Connect(ctx context.Context) (err error) {
	_, spanFinish := c.startTraceSpan(&ctx, "connect", &err, nil)
	defer spanFinish()
	return c.client.Connect(ctx)
}

// Disconnect is a wrapper.
func (c *Client) Disconnect(ctx context.Context) (err error) {
	_, spanFinish := c.startTraceSpan(&ctx, "disconnect", &err, nil)
	defer spanFinish()
	return c.client.Disconnect(ctx)
}

// Ping is a wrapper.
func (c *Client) Ping(ctx context.Context, rp *readpref.ReadPref) (err error) {
	_, spanFinish := c.startTraceSpan(&ctx, "ping", &err, nil)
	defer spanFinish()
	return c.client.Ping(ctx, rp)
}

// Database is a wrapper.
func (c *Client) Database(name string, opts ...*options.DatabaseOptions) *Database {
	return &Database{
		database: c.client.Database(name, opts...),
		client:   c,
	}
}

// ListDatabaseNames is a wrapper.
func (c *Client) ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (dbNames []string, err error) {
	span, spanFinish := c.startTraceSpan(&ctx, "list_database_names", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	return c.client.ListDatabaseNames(ctx, filter, opts...)
}

// ListDatabases is a wrapper.
func (c *Client) ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (res ListDatabasesResult, err error) {
	span, spanFinish := c.startTraceSpan(&ctx, "list_databases", &err, func(err error) error {
		if filter != nil {
			err = wrapErrorValue(err, "filter", filter)
		}
		return err
	})
	defer spanFinish()
	if filter != nil {
		setTraceSpanTagJSON(span, "filter", filter)
	}
	return c.client.ListDatabases(ctx, filter, opts...)
}

// StartSession is a wrapper.
func (c *Client) StartSession(opts ...*options.SessionOptions) (s Session, err error) {
	defer c.wrapErrorReturn("start_session", &err, nil)
	return c.client.StartSession(opts...)
}

// UseSession is a wrapper.
func (c *Client) UseSession(ctx context.Context, fn func(SessionContext) error) (err error) {
	defer c.wrapErrorReturn("use_session", &err, nil)
	return c.client.UseSession(ctx, fn)
}

// UseSessionWithOptions is a wrapper.
func (c *Client) UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(SessionContext) error) (err error) {
	defer c.wrapErrorReturn("use_session_with_options", &err, nil)
	return c.client.UseSessionWithOptions(ctx, opts, fn)
}

// Watch is a wrapper.
func (c *Client) Watch(ctx context.Context, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (_ *ChangeStream, err error) {
	defer c.wrapErrorReturn("watch", &err, func(err error) error {
		if pipeline != nil {
			err = wrapErrorValue(err, "pipeline", pipeline)
		}
		return err
	})
	cs, err := c.client.Watch(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}
	return &ChangeStream{
		changeStream: cs,
		client:       c,
	}, nil
}

func (c *Client) startTraceSpan(pctx *context.Context, op string, perr *error, werr func(error) error) (opentracing.Span, closeutils.F) {
	span, spanFinish := startTraceSpan(pctx, "client."+op, perr, func(err error) error {
		err = c.wrapErrorValues(err)
		err = wrapErrorOptional(err, werr)
		return err
	})
	c.setTraceSpanTags(span)
	return span, spanFinish
}

func (c *Client) setTraceSpanTags(span opentracing.Span) {
	setClientOptionsTraceSpanTags(c.options, span)
}

func (c *Client) wrapErrorReturn(op string, perr *error, werr func(error) error) {
	wrapErrorReturn("client."+op, perr, func(err error) error {
		err = c.wrapErrorValues(err)
		err = wrapErrorOptional(err, werr)
		return err
	})
}

func (c *Client) wrapErrorValues(err error) error {
	return wrapClientOptionsErrorValues(c.options, err)
}

func setClientOptionsTraceSpanTags(opts *options.ClientOptions, span opentracing.Span) {
	if len(opts.Hosts) != 0 {
		setTraceSpanTag(span, "client_options.hosts", strings.Join(opts.Hosts, " "))
	}
	if opts.ReplicaSet != nil {
		setTraceSpanTag(span, "client_options.replica_set", *opts.ReplicaSet)
	}
	if opts.AppName != nil {
		setTraceSpanTag(span, "client_options.app_name", *opts.AppName)
	}
	if opts.ConnectTimeout != nil {
		setTraceSpanTag(span, "client_options.connect_timeout", opts.ConnectTimeout.String())
	}
	if opts.SocketTimeout != nil {
		setTraceSpanTag(span, "client_options.socket_timeout", opts.SocketTimeout.String())
	}
	if opts.ServerSelectionTimeout != nil {
		setTraceSpanTag(span, "client_options.server_selection_timeout", opts.ServerSelectionTimeout.String())
	}
}

func wrapClientOptionsErrorValues(opts *options.ClientOptions, err error) error {
	if len(opts.Hosts) != 0 {
		err = wrapErrorValue(err, "client_options.hosts", opts.Hosts)
	}
	if opts.ReplicaSet != nil {
		err = wrapErrorValue(err, "client_options.replica_set", *opts.ReplicaSet)
	}
	if opts.AppName != nil {
		err = wrapErrorValue(err, "client_options.app_name", *opts.AppName)
	}
	if opts.ConnectTimeout != nil {
		err = wrapErrorValue(err, "client_options.connect_timeout", *opts.ConnectTimeout)
	}
	if opts.SocketTimeout != nil {
		err = wrapErrorValue(err, "client_options.socket_timeout", *opts.SocketTimeout)
	}
	if opts.ServerSelectionTimeout != nil {
		err = wrapErrorValue(err, "client_options.server_selection_timeout", *opts.ServerSelectionTimeout)
	}
	return err
}
