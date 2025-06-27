package nats

import (
	"fmt"
	"time"

	natsio "github.com/nats-io/nats.go"
)

const (
	_defaultConnTimeout = 3 * time.Second
)

// NatsClient wraps the NATS connection and provides basic pub/sub methods.
type NatsClient struct {
	conn    *natsio.Conn
	timeout time.Duration
}

// New creates a new NatsClient.
func New(url string, opts ...Option) (*NatsClient, error) {
	nc := &NatsClient{
		timeout: _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(nc)
	}

	conn, err := natsio.Connect(url, natsio.Timeout(nc.timeout))
	if err != nil {
		return nil, fmt.Errorf("nats - New - nats.Connect: %w", err)
	}

	nc.conn = conn
	return nc, nil
}

// Publish sends a message to a subject.
func (n *NatsClient) Publish(subject string, data []byte) error {
	return n.conn.Publish(subject, data)
}

// Subscribe subscribes to a subject and handles messages with the given handler.
func (n *NatsClient) Subscribe(subject string, handler func(msg *natsio.Msg)) (*natsio.Subscription, error) {
	return n.conn.Subscribe(subject, handler)
}

// Close closes the NATS connection.
func (n *NatsClient) Close() error {
	if n.conn != nil {
		n.conn.Close()
	}
	return nil
}

// Option for configuring NatsClient.
type Option func(*NatsClient)

// ConnTimeout sets the connection timeout.
func ConnTimeout(timeout time.Duration) Option {
	return func(n *NatsClient) {
		n.timeout = timeout
	}
}
