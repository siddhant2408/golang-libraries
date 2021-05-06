package mongoutils

import (
	"context"
	"net/url"
	"strings"
	"sync"

	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/goroutine"
	"github.com/siddhant2408/golang-libraries/mongo"
	"github.com/siddhant2408/golang-libraries/tracingutils"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ClientPool is a pool of MongoDB clients.
type ClientPool struct {
	mu  sync.RWMutex
	its map[string]*clientPoolItem
	// Connect allows to customize the function that connects to MongoDB.
	// By default it uses mongo.Connect.
	Connect func(context.Context, ...*options.ClientOptions) (*mongo.Client, error)
}

// NewClientPool returns a new ClientPool.
func NewClientPool() *ClientPool {
	return &ClientPool{
		its: make(map[string]*clientPoolItem),
	}
}

// GetClient returns a client from the pool for the given URI.
func (cp *ClientPool) GetClient(ctx context.Context, u string) (*mongo.Client, error) {
	return cp.getItem(u).getClient(ctx)
}

func (cp *ClientPool) getItem(u string) *clientPoolItem {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	it, ok := cp.its[u]
	if ok {
		return it
	}
	it = &clientPoolItem{
		uri:     u,
		connect: cp.Connect,
	}
	cp.its[u] = it
	return it
}

// Close closes the ClientPool.
func (cp *ClientPool) Close(ctx context.Context, oe closeutils.OnErr) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	wg := new(sync.WaitGroup)
	for u, it := range cp.its {
		u, it := u, it
		goroutine.WaitGroup(wg, func() {
			err := it.close(ctx)
			if err != nil {
				err = errors.WithValue(err, "mongo.uri", u)
				oe(err)
			}
		})
		delete(cp.its, u)
	}
	wg.Wait()
}

// ForceClose closes the ClientPool with a context that is already canceled.
// It forces closes all opened connections.
func (cp *ClientPool) ForceClose(oe closeutils.OnErr) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	cp.Close(ctx, oe)
}

type clientPoolItem struct {
	mu      sync.Mutex
	clt     *mongo.Client
	uri     string
	connect func(context.Context, ...*options.ClientOptions) (*mongo.Client, error)
}

func (it *clientPoolItem) getClient(ctx context.Context) (clt *mongo.Client, err error) {
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "mongoutils.client_pool", &err)
	defer spanFinish()
	span.SetTag("mongo.uri", it.uri)
	tracingutils.TraceSyncLocker(ctx, &it.mu)
	defer it.mu.Unlock()
	if it.clt != nil {
		return it.clt, nil
	}
	opts := options.Client().ApplyURI(it.uri)
	connect := it.connect
	if connect == nil {
		connect = mongo.Connect
	}
	clt, err = connect(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	it.clt = clt
	return it.clt, nil
}

func (it *clientPoolItem) close(ctx context.Context) error {
	it.mu.Lock()
	defer it.mu.Unlock()
	if it.clt == nil {
		return nil
	}
	clt := it.clt
	it.clt = nil
	err := clt.Disconnect(ctx)
	if err != nil {
		return errors.Wrap(err, "disconnect")
	}
	return nil
}

// DatabasePool is a pool of MongoDB databases.
//
// It allows to reuse MongoDB client if the base connection string (without DB name) is the same.
type DatabasePool struct {
	// ClientPool is the pool of clients.
	ClientPool interface {
		GetClient(ctx context.Context, u string) (*mongo.Client, error)
	}
	// Database is a function that allows to customize the function that get the dabatabase from the client.
	// By default it uses mongo.Client.Database.
	Database func(clt *mongo.Client, name string, opts ...*options.DatabaseOptions) *mongo.Database
}

// GetDatabase returns a database for the given URI.
func (dp *DatabasePool) GetDatabase(ctx context.Context, u string) (db *mongo.Database, err error) {
	span, spanFinish := tracingutils.StartChildSpan(&ctx, "mongoutils.database_pool", &err)
	defer spanFinish()
	span.SetTag("mongo.uri", u)
	baseURI, dbName, err := dp.parseURI(u)
	if err != nil {
		return nil, errors.Wrap(err, "parse URI")
	}
	clt, err := dp.ClientPool.GetClient(ctx, baseURI)
	if err != nil {
		return nil, errors.Wrap(err, "get client")
	}
	database := dp.Database
	if database == nil {
		database = (*mongo.Client).Database
	}
	db = database(clt, dbName)
	return db, nil
}

func (dp *DatabasePool) parseURI(u string) (baseURI string, dbName string, err error) {
	pu, err := url.Parse(u)
	if err != nil {
		return "", "", errors.Wrap(err, "")
	}
	dbName = strings.TrimPrefix(pu.Path, "/")
	if dbName == "" {
		return "", "", errors.New("empty database name")
	}
	pu.Path = ""
	baseURI = pu.String()
	return baseURI, dbName, nil
}
