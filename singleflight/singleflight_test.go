package singleflight

import (
	"context"
	"testing"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestDo(t *testing.T) {
	ctx := context.Background()
	g := &Group{}
	v, err, _ := g.Do(ctx, "test", func(ctx context.Context) (interface{}, error) {
		return "test", nil
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	expected := "test"
	if v != expected {
		t.Fatalf("unexpected value: got %v, want %q", v, expected)
	}
}

func TestDoError(t *testing.T) {
	ctx := context.Background()
	g := &Group{}
	_, err, _ := g.Do(ctx, "test", func(ctx context.Context) (interface{}, error) {
		return nil, errors.New("error")
	})
	if err == nil {
		t.Fatal("no error")
	}
}
