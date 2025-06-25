// Package redis implements a Redis client.
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	_defaultConnTimeout = 3 * time.Second
)

// Redis -.
type Redis struct {
	client  *redis.Client
	timeout time.Duration
}

// New -.
func New(url string, opts ...Option) (*Redis, error) {
	r := &Redis{
		timeout: _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(r)
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("redis - New - redis.ParseURL: %w", err)
	}

	r.client = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	if err := r.client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis - New - r.client.Ping: %w", err)
	}

	return r, nil
}

// Client -.
func (r *Redis) Client() *redis.Client {
	return r.client
}

// Close -.
func (r *Redis) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}
