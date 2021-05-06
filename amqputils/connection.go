package amqputils

import (
	"context"

	"github.com/siddhant2408/golang-libraries/ctxsync"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
	"github.com/streadway/amqp"
)

// ConnectionManager manages a amqp.Connection.
type ConnectionManager struct {
	Dial func(context.Context) (*amqp.Connection, error)

	mu   ctxsync.Mutex
	conn *amqp.Connection
}

// Channel opens a new amqp.Channel on the managed amqp.Connection.
//
// If an error occurs during the amqp.Channel opening, the ConnectionManager is closed.
func (m *ConnectionManager) Channel(ctx context.Context) (chn *amqp.Channel, err error) {
	_, spanFinish := startTraceChildSpan(&ctx, "connection_manager.channel", &err)
	defer spanFinish()
	err = tracingutils.TraceSyncLockerCtx(ctx, &m.mu)
	if err != nil {
		return nil, errors.Wrap(err, "lock")
	}
	defer m.mu.Unlock()
	conn, err := m.get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get connection")
	}
	chn, err = connectionChannel(ctx, conn)
	if err != nil {
		_ = m.close()
		return nil, errors.Wrap(err, "open channel")
	}
	return chn, nil
}

func (m *ConnectionManager) get(ctx context.Context) (conn *amqp.Connection, err error) {
	if m.conn != nil {
		return m.conn, nil
	}
	conn, err = m.Dial(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "dial")
	}
	m.conn = conn
	return conn, nil
}

// Close closes the managed amqp.Connection and unsets it.
// It does nothing if the amqp.Connection is not set.
//
// It is OK to reuse the ConnectionManager after this call.
func (m *ConnectionManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.close()
}

func (m *ConnectionManager) close() error {
	if m.conn == nil {
		return nil
	}
	err := m.conn.Close()
	m.conn = nil
	return errors.Wrap(err, "close connection")
}

// NewConnectionManagerURLs returns a new ConnectionManager for the given URLs.
func NewConnectionManagerURLs(urls []string) *ConnectionManager {
	d := &URLsDialer{
		URLs:    urls,
		DialURL: Dial,
	}
	return &ConnectionManager{
		Dial: d.Dial,
	}
}
