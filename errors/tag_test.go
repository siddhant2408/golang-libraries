package errors_test

import (
	"fmt"
	"io"
	"testing"

	. "github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/errors/internal"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestTag(t *testing.T) {
	err := internal.NewBase("error")
	err = WithTag(err, "foo", "bar")
	tags := Tags(err)
	expected := map[string]string{
		"foo": "bar",
	}
	testutils.Compare(t, "unexpected tags", tags, expected)
}

func TestTagInt(t *testing.T) {
	err := internal.NewBase("error")
	err = WithTagInt(err, "foo", 123)
	tags := Tags(err)
	expected := map[string]string{
		"foo": "123",
	}
	testutils.Compare(t, "unexpected tags", tags, expected)
}

func TestTagInt64(t *testing.T) {
	err := internal.NewBase("error")
	err = WithTagInt64(err, "foo", 123)
	tags := Tags(err)
	expected := map[string]string{
		"foo": "123",
	}
	testutils.Compare(t, "unexpected tags", tags, expected)
}

func TestTagFloat64(t *testing.T) {
	err := internal.NewBase("error")
	err = WithTagFloat64(err, "foo", 12.3)
	tags := Tags(err)
	expected := map[string]string{
		"foo": "12.3",
	}
	testutils.Compare(t, "unexpected tags", tags, expected)
}

func TestTagBool(t *testing.T) {
	err := internal.NewBase("error")
	err = WithTagBool(err, "foo", true)
	tags := Tags(err)
	expected := map[string]string{
		"foo": "true",
	}
	testutils.Compare(t, "unexpected tags", tags, expected)
}

func TestTagOverWrite(t *testing.T) {
	err := internal.NewBase("error")
	err = WithTag(err, "test", "1")
	err = WithTag(err, "test", "2")
	tags := Tags(err)
	expected := map[string]string{
		"test": "2",
	}
	testutils.Compare(t, "unexpected tags", tags, expected)
}

func TestTagNil(t *testing.T) {
	err := WithTag(nil, "foo", "bar")
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestTagsEmpty(t *testing.T) {
	err := internal.NewBase("error")
	tags := Tags(err)
	if len(tags) != 0 {
		t.Fatalf("tags not empty: got %#v", tags)
	}
}

func TestTagError(t *testing.T) {
	err := internal.NewBase("error")
	err = WithTag(err, "foo", "bar")
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestTagFormat(t *testing.T) {
	err := internal.NewBase("error")
	err = WithTag(err, "foo", "bar")
	s := fmt.Sprintf("%+v", err)
	expected := "tag foo = bar\nerror"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func BenchmarkTagFormat(b *testing.B) {
	err := internal.NewBase("error")
	err = WithTag(err, "foo", "bar")
	for i := 0; i < b.N; i++ {
		_, _ = fmt.Fprintf(io.Discard, "%+v", err)
	}
}
