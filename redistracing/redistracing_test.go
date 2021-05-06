package redistracing

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/redistest"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestHookSet(t *testing.T) {
	ctx := context.Background()
	clt := redistest.NewClient(t)
	opts := clt.Options()
	clt.AddHook(&Hook{
		Addrs: []string{opts.Addr},
		DB:    &opts.DB,
	})
	err := clt.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestHookDel(t *testing.T) {
	ctx := context.Background()
	clt := redistest.NewClient(t)
	opts := clt.Options()
	clt.AddHook(&Hook{
		Addrs: []string{opts.Addr},
		DB:    &opts.DB,
	})
	err := clt.Del(ctx, "foo").Err()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestHookMSet(t *testing.T) {
	ctx := context.Background()
	clt := redistest.NewClient(t)
	opts := clt.Options()
	clt.AddHook(&Hook{
		Addrs: []string{opts.Addr},
		DB:    &opts.DB,
	})
	err := clt.MSet(ctx, "foo", "bar").Err()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestHookBRPop(t *testing.T) {
	ctx := context.Background()
	clt := redistest.NewClient(t)
	opts := clt.Options()
	clt.AddHook(&Hook{
		Addrs: []string{opts.Addr},
		DB:    &opts.DB,
	})
	err := clt.LPush(ctx, "foo", "bar").Err()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = clt.BRPop(ctx, 0, "foo").Err()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestHookPipeline(t *testing.T) {
	ctx := context.Background()
	clt := redistest.NewClient(t)
	opts := clt.Options()
	clt.AddHook(&Hook{
		Addrs: []string{opts.Addr},
		DB:    &opts.DB,
	})
	p := clt.TxPipeline()
	defer p.Close() //nolint:errcheck
	p.Set(ctx, "foo", "bar", 0)
	_, err := p.Exec(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
