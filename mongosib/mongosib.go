// Package mongosib provides MongoDB related utilities.
package mongosib

import (
	"context"
	"time"

	"github.com/siddhant2408/golang-libraries/envutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/mongo"
	"github.com/siddhant2408/golang-libraries/mongosib/mongosibbson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var clientOptionsDefault = options.Client().
	SetConnectTimeout(10 * time.Second).
	SetServerSelectionTimeout(30 * time.Second).
	SetSocketTimeout(1 * time.Minute).
	SetMaxConnIdleTime(1 * time.Minute).
	SetReadPreference(readpref.PrimaryPreferred()).
	SetRegistry(mongosibbson.NewRegistry())

// Connect connects with the default options.
func Connect(ctx context.Context, appName string, opts ...*options.ClientOptions) (*mongo.Client, error) {
	if appName == "" {
		panic("application name must not be empty")
	}
	opts = append([]*options.ClientOptions{clientOptionsDefault}, opts...)
	opts = append(opts, options.Client().SetAppName(appName))
	return mongo.Connect(ctx, opts...)
}

// ConnectURIsEnv connects with a list of URIs per environments.
func ConnectURIsEnv(ctx context.Context, uris map[envutils.Env]string, env envutils.Env, appName string, opts ...*options.ClientOptions) (*mongo.Client, error) {
	u, err := getURIEnv(uris, env)
	if err != nil {
		return nil, errors.Wrap(err, "get URI")
	}
	opts = append([]*options.ClientOptions{options.Client().ApplyURI(u)}, opts...)
	return Connect(ctx, appName, opts...)
}

func getURIEnv(uris map[envutils.Env]string, env envutils.Env) (string, error) {
	u, ok := uris[env]
	if !ok {
		return "", errors.Newf("invalid environment %q", env)
	}
	return u, nil
}

var urisCentralDB = map[envutils.Env]string{
	envutils.Testing:     "mongodb://localhost:27017",
	envutils.Development: "",
	envutils.Staging:     "",
	envutils.Production:  "",
}

// ConnectCentralDB connects to central DB.
func ConnectCentralDB(ctx context.Context, env envutils.Env, appName string, opts ...*options.ClientOptions) (*mongo.Client, error) {
	return ConnectURIsEnv(ctx, urisCentralDB, env, appName, opts...)
}

var urisMTAProcess = map[envutils.Env]string{
	envutils.Testing:     "mongodb://localhost:27017",
	envutils.Development: "",
	envutils.Staging:     "",
	envutils.Production:  "",
}

// ConnectMTAProcess connects to MTA process.
func ConnectMTAProcess(ctx context.Context, env envutils.Env, appName string, opts ...*options.ClientOptions) (*mongo.Client, error) {
	return ConnectURIsEnv(ctx, urisMTAProcess, env, appName, opts...)
}
