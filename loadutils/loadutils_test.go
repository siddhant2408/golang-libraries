package loadutils

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestLoader(t *testing.T) {
	k := "test"
	v1 := "value"
	cg := &testCacheGetter{
		found: false,
	}
	cs := &testCacheSetter{}
	c := &testCache{
		testCacheSetter: cs,
		testCacheGetter: cg,
	}
	l := &Loader{
		Cache: c,
	}
	fCalled := false
	f := func(ctx context.Context) (value interface{}, cache bool, err error) {
		fCalled = true
		return v1, true, nil
	}
	v2i, err := l.Load(context.Background(), k, f)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if v2i != v1 {
		t.Fatalf("unexpected value: got %v, want %q", v2i, v1)
	}
	cg.checkCalled(t)
	cs.checkCalled(t)
	if !fCalled {
		t.Fatal("not called")
	}
}

func TestLoaderCacheGetFound(t *testing.T) {
	k := "test"
	v1 := "value"
	cg := &testCacheGetter{
		value: v1,
		found: true,
	}
	c := &testCache{
		testCacheGetter: cg,
	}
	l := &Loader{
		Cache: c,
	}
	v2i, err := l.Load(context.Background(), k, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if v2i != v1 {
		t.Fatalf("unexpected value: got %v, want %q", v2i, v1)
	}
	cg.checkCalled(t)
}

func TestLoaderNoCache(t *testing.T) {
	k := "test"
	v1 := "value"
	cg := &testCacheGetter{
		found: false,
	}
	c := &testCache{
		testCacheGetter: cg,
	}
	l := &Loader{
		Cache: c,
	}
	fCalled := false
	f := func(ctx context.Context) (value interface{}, cache bool, err error) {
		fCalled = true
		return v1, false, nil
	}
	v2i, err := l.Load(context.Background(), k, f)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if v2i != v1 {
		t.Fatalf("unexpected value: got %v, want %q", v2i, v1)
	}
	cg.checkCalled(t)
	if !fCalled {
		t.Fatal("not called")
	}
}

func TestLoaderErrorCacheGet(t *testing.T) {
	k := "test"
	cg := &testCacheGetter{
		err: errors.New("error"),
	}
	c := &testCache{
		testCacheGetter: cg,
	}
	l := &Loader{
		Cache: c,
	}
	_, err := l.Load(context.Background(), k, nil)
	if err == nil {
		t.Fatal("no error")
	}
	cg.checkCalled(t)
}

func TestLoaderErrorLoad(t *testing.T) {
	k := "test"
	cg := &testCacheGetter{
		found: false,
	}
	c := &testCache{
		testCacheGetter: cg,
	}
	l := &Loader{
		Cache: c,
	}
	fCalled := false
	f := func(ctx context.Context) (value interface{}, cache bool, err error) {
		fCalled = true
		return nil, false, errors.New("error")
	}
	_, err := l.Load(context.Background(), k, f)
	if err == nil {
		t.Fatal("no error")
	}
	cg.checkCalled(t)
	if !fCalled {
		t.Fatal("not called")
	}
}

func TestLoaderErrorCacheSet(t *testing.T) {
	k := "test"
	v1 := "value"
	cg := &testCacheGetter{
		found: false,
	}
	cs := &testCacheSetter{
		err: errors.New("error"),
	}
	c := &testCache{
		testCacheSetter: cs,
		testCacheGetter: cg,
	}
	l := &Loader{
		Cache: c,
	}
	fCalled := false
	f := func(ctx context.Context) (value interface{}, cache bool, err error) {
		fCalled = true
		return v1, true, nil
	}
	_, err := l.Load(context.Background(), k, f)
	if err == nil {
		t.Fatal("no error")
	}
	cg.checkCalled(t)
	cs.checkCalled(t)
	if !fCalled {
		t.Fatal("not called")
	}
}

type testCache struct {
	*testCacheSetter
	*testCacheGetter
}

type testCacheSetter struct {
	called bool
	err    error
}

func (c *testCacheSetter) Set(_ context.Context, _ string, _ interface{}) error {
	c.called = true
	return c.err
}

func (c *testCacheSetter) checkCalled(tb testing.TB) {
	tb.Helper()
	if !c.called {
		tb.Fatal("not called")
	}
}

type testCacheGetter struct {
	called bool
	value  interface{}
	found  bool
	err    error
}

func (c *testCacheGetter) Get(_ context.Context, _ string) (value interface{}, found bool, err error) {
	c.called = true
	return c.value, c.found, c.err
}

func (c *testCacheGetter) checkCalled(tb testing.TB) {
	tb.Helper()
	if !c.called {
		tb.Fatal("not called")
	}
}

func TestInMemoryCacheWrapper(t *testing.T) {
	w := &InMemoryCacheWrapper{
		Cache: &testInMemoryCache{},
	}
	err := w.Set(context.Background(), "key", "value")
	if err != nil {
		testutils.FatalErr(t, err)
	}
	_, found, err := w.Get(context.Background(), "key")
	if err != nil {
		testutils.FatalErr(t, err)
	}
	if !found {
		t.Fatal("not found")
	}
}

type testInMemoryCache struct{}

func (c *testInMemoryCache) Set(key interface{}, value interface{}) {
}

func (c *testInMemoryCache) Get(key interface{}) (value interface{}, found bool) {
	return "value", true
}
