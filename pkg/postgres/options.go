package postgres

import "time"

// Option -.
type Option func(*Postgres)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

// MinPoolSize -.
func MinPoolSize(size int) Option {
	return func(c *Postgres) {
		c.minPoolSize = size
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *Postgres) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.connTimeout = timeout
	}
}

// MaxConnLifetime -.
func MaxConnLifetime(lifetime time.Duration) Option {
	return func(c *Postgres) {
		c.maxConnLifetime = lifetime
	}
}

// MaxConnIdleTime -.
func MaxConnIdleTime(idleTime time.Duration) Option {
	return func(c *Postgres) {
		c.maxConnIdleTime = idleTime
	}
}

// HealthCheckPeriod -.
func HealthCheckPeriod(period time.Duration) Option {
	return func(c *Postgres) {
		c.healthCheckPeriod = period
	}
}
