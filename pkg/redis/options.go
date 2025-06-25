package redis

import "time"

// Option -.
type Option func(*Redis)

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Redis) {
		c.timeout = timeout
	}
}
