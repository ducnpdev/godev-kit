// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxPoolSize       = 100
	_defaultMinPoolSize       = 10
	_defaultConnAttempts      = 10
	_defaultConnTimeout       = time.Second
	_defaultMaxConnLifetime   = time.Hour
	_defaultMaxConnIdleTime   = time.Minute * 30
	_defaultHealthCheckPeriod = time.Minute
)

// Postgres -.
type Postgres struct {
	maxPoolSize       int
	minPoolSize       int
	connAttempts      int
	connTimeout       time.Duration
	maxConnLifetime   time.Duration
	maxConnIdleTime   time.Duration
	healthCheckPeriod time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

// New -.
func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:       _defaultMaxPoolSize,
		minPoolSize:       _defaultMinPoolSize,
		connAttempts:      _defaultConnAttempts,
		connTimeout:       _defaultConnTimeout,
		maxConnLifetime:   _defaultMaxConnLifetime,
		maxConnIdleTime:   _defaultMaxConnIdleTime,
		healthCheckPeriod: _defaultHealthCheckPeriod,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize) //nolint:gosec // skip integer overflow conversion int -> int32
	poolConfig.MinConns = int32(pg.minPoolSize) //nolint:gosec // skip integer overflow conversion int -> int32
	poolConfig.MaxConnLifetime = pg.maxConnLifetime
	poolConfig.MaxConnIdleTime = pg.maxConnIdleTime
	poolConfig.HealthCheckPeriod = pg.healthCheckPeriod

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			err = pg.Ping(context.Background())
			if err != nil {
				return nil, fmt.Errorf("postgres - NewPostgres - pg.Ping: %w", err)
			}
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

// Close -.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

// Ping tests the database connection
func (p *Postgres) Ping(ctx context.Context) error {
	if p.Pool == nil {
		return fmt.Errorf("postgres pool is not initialized")
	}

	return p.Pool.Ping(ctx)
}

// GetPoolStats returns connection pool statistics for monitoring
func (p *Postgres) GetPoolStats() map[string]interface{} {
	if p.Pool == nil {
		return map[string]interface{}{
			"total_connections":  0,
			"idle_connections":   0,
			"in_use_connections": 0,
			"max_connections":    0,
		}
	}

	stats := p.Pool.Stat()
	return map[string]interface{}{
		"total_connections":  stats.TotalConns(),
		"idle_connections":   stats.IdleConns(),
		"in_use_connections": stats.TotalConns() - stats.IdleConns(),
		"max_connections":    stats.MaxConns(),
	}
}
