// Package jsonlog provides a helper to write JSON log.
package jsonlog

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

// Logger represents a JSON logger.
// It is safe to use it concurrently.
type Logger struct {
	mu  sync.Mutex
	enc *json.Encoder
}

// New returns a new Logger for a io.Writer.
func New(w io.Writer) *Logger {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return &Logger{
		enc: enc,
	}
}

// NewFile create a new Logger that writes to a file.
func NewFile(name string, perm os.FileMode) (*Logger, closeutils.Err, error) {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, perm)
	if err != nil {
		return nil, nil, errors.Wrap(err, "open file")
	}
	cl := func() error {
		err = f.Close()
		if err != nil {
			return errors.Wrap(err, "close file")
		}
		return nil
	}
	return New(f), cl, nil
}

// Log writes the JSON log.
func (l *Logger) Log(ctx context.Context, data interface{}) (err error) {
	_, spanFinish := tracingutils.StartChildSpan(&ctx, "jsonlog", &err)
	defer spanFinish()
	l.mu.Lock()
	defer l.mu.Unlock()
	err = l.enc.Encode(data)
	if err != nil {
		return errors.Wrap(err, "encode")
	}
	return nil
}

// Optional is a logger that optionally logs if a logger is set.
type Optional struct {
	Logger interface {
		Log(context.Context, interface{}) error
	}
}

// NewOptionalFile returns a new Optional that writes to a file.
// If the name is empty, the Optional doesn't write anything.
func NewOptionalFile(name string, perm os.FileMode) (*Optional, closeutils.Err, error) {
	if name == "" {
		return &Optional{}, func() error { return nil }, nil
	}
	l, cl, err := NewFile(name, perm)
	if err != nil {
		return nil, nil, errors.Wrap(err, "new")
	}
	return &Optional{
		Logger: l,
	}, cl, nil
}

// Log optionally writes the JSON log.
func (l *Optional) Log(ctx context.Context, data interface{}) error {
	if l.Logger != nil {
		return l.Logger.Log(ctx, data)
	}
	return nil
}

// Error is a logger that handles the error instead of returning it.
//
// If Logger returns an error, it calls OnError.
type Error struct {
	Logger interface {
		Log(context.Context, interface{}) error
	}
	OnError func(context.Context, error)
}

// NewErrorOptionalFile returns a new Error that writes to an optional file.
func NewErrorOptionalFile(name string, perm os.FileMode, onError func(context.Context, error)) (*Error, closeutils.Err, error) {
	l, cl, err := NewOptionalFile(name, perm)
	if err != nil {
		return nil, nil, errors.Wrap(err, "optional")
	}
	return &Error{
		Logger:  l,
		OnError: onError,
	}, cl, nil
}

// Log optionally writes the JSON log.
func (l *Error) Log(ctx context.Context, data interface{}) {
	err := l.Logger.Log(ctx, data)
	if err != nil {
		err = errors.Wrap(err, "jsonlog")
		l.OnError(ctx, err)
	}
}
