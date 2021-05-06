package mongosib

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/envutils"
	"github.com/siddhant2408/golang-libraries/mongotest"
	"github.com/siddhant2408/golang-libraries/testutils"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestConnect(t *testing.T) {
	mongotest.CheckAvailable(t)
	ctx := context.Background()
	clt, err := Connect(ctx, "test")
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = clt.Ping(ctx, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = clt.Disconnect(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestConnectError(t *testing.T) {
	ctx := context.Background()
	opts := options.Client().ApplyURI("invalid")
	_, err := Connect(ctx, "test", opts)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestConnectPanicAppNameEmpty(t *testing.T) {
	ctx := context.Background()
	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatal("no panic")
		}
	}()
	_, _ = Connect(ctx, "")
}

func TestConnectURIsEnv(t *testing.T) {
	mongotest.CheckAvailable(t)
	ctx := context.Background()
	uris := map[envutils.Env]string{
		envutils.Testing: mongotest.GetURI(),
	}
	clt, err := ConnectURIsEnv(ctx, uris, envutils.Testing, "test")
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = clt.Ping(ctx, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = clt.Disconnect(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestConnectURIsEnvErrorEnv(t *testing.T) {
	ctx := context.Background()
	uris := map[envutils.Env]string{
		envutils.Testing: mongotest.GetURI(),
	}
	_, err := ConnectURIsEnv(ctx, uris, envutils.Env("invalid"), "test")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestConnectCentralDB(t *testing.T) {
	mongotest.CheckAvailable(t)
	ctx := context.Background()
	clt, err := ConnectCentralDB(ctx, envutils.Testing, "test")
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = clt.Ping(ctx, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = clt.Disconnect(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestConnectMTAProcess(t *testing.T) {
	mongotest.CheckAvailable(t)
	ctx := context.Background()
	clt, err := ConnectMTAProcess(ctx, envutils.Testing, "test")
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = clt.Ping(ctx, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = clt.Disconnect(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
