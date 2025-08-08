package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ducnpdev/godev-kit/pkg/postgres"
)

func main() {
	// Example 1: Basic configuration with default values
	fmt.Println("=== Example 1: Basic Configuration ===")
	pg1, err := postgres.New("postgres://username:password@localhost:5432/database")
	if err != nil {
		log.Fatal(err)
	}
	defer pg1.Close()
	fmt.Printf("Postgres connected with default settings\n")

	// Example 2: Custom pool configuration
	fmt.Println("\n=== Example 2: Custom Pool Configuration ===")
	pg2, err := postgres.New(
		"postgres://username:password@localhost:5432/database",
		postgres.MaxPoolSize(20),
		postgres.MinPoolSize(5),
		postgres.MaxConnLifetime(time.Hour),
		postgres.MaxConnIdleTime(30*time.Minute),
		postgres.HealthCheckPeriod(time.Minute),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer pg2.Close()
	fmt.Printf("Postgres connected with custom pool settings\n")

	// Example 3: High-performance configuration
	fmt.Println("\n=== Example 3: High-Performance Configuration ===")
	pg3, err := postgres.New(
		"postgres://username:password@localhost:5432/database",
		postgres.MaxPoolSize(50),
		postgres.MinPoolSize(10),
		postgres.MaxConnLifetime(30*time.Minute),
		postgres.MaxConnIdleTime(10*time.Minute),
		postgres.HealthCheckPeriod(15*time.Second),
		postgres.ConnAttempts(5),
		postgres.ConnTimeout(5*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer pg3.Close()
	fmt.Printf("Postgres connected with high-performance settings\n")

	// Example 4: Monitoring pool statistics
	fmt.Println("\n=== Example 4: Pool Statistics ===")
	stats := pg3.Pool.Stat()
	fmt.Printf("Total connections: %d\n", stats.TotalConns())
	fmt.Printf("Idle connections: %d\n", stats.IdleConns())
	fmt.Printf("In-use connections: %d\n", stats.InUseConns())
	fmt.Printf("Max connections: %d\n", stats.MaxConns())

	// Example 5: Simple query to test connection
	fmt.Println("\n=== Example 5: Testing Connection ===")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result string
	err = pg3.Pool.QueryRow(ctx, "SELECT version()").Scan(&result)
	if err != nil {
		log.Printf("Query failed: %v", err)
	} else {
		fmt.Printf("Database version: %s\n", result)
	}

	fmt.Println("\n=== Connection Pool Demo Completed ===")
}

// Example configuration for different environments
func getEnvironmentConfig(env string) []postgres.Option {
	switch env {
	case "development":
		return []postgres.Option{
			postgres.MaxPoolSize(5),
			postgres.MinPoolSize(1),
			postgres.MaxConnLifetime(time.Hour),
			postgres.MaxConnIdleTime(30*time.Minute),
			postgres.HealthCheckPeriod(time.Minute),
		}
	case "production":
		return []postgres.Option{
			postgres.MaxPoolSize(20),
			postgres.MinPoolSize(5),
			postgres.MaxConnLifetime(time.Hour),
			postgres.MaxConnIdleTime(15*time.Minute),
			postgres.HealthCheckPeriod(30*time.Second),
		}
	case "high-traffic":
		return []postgres.Option{
			postgres.MaxPoolSize(50),
			postgres.MinPoolSize(10),
			postgres.MaxConnLifetime(30*time.Minute),
			postgres.MaxConnIdleTime(10*time.Minute),
			postgres.HealthCheckPeriod(15*time.Second),
		}
	default:
		return []postgres.Option{}
	}
} 