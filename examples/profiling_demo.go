package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ducnpdev/godev-kit/pkg/profiling"
	"github.com/rs/zerolog"
)

// Example of how to use the profiling system

func profilingDemo() {
	// Initialize logger
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Initialize profiler
	profiler := profiling.NewProfiler(logger, true, "/debug")

	// Example 1: Profile a simple function
	fmt.Println("=== Example 1: Simple Function Profiling ===")

	err := profiler.ProfileFunction("example_function", "main", func() error {
		// Simulate some work
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	if err != nil {
		log.Printf("Error: %v", err)
	}

	// Example 2: Profile a function with context
	fmt.Println("\n=== Example 2: Function with Context Profiling ===")

	ctx := context.Background()
	err = profiler.ProfileFunctionWithContext(ctx, "example_context_function", "main", func(ctx context.Context) error {
		// Simulate some work with context
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(50 * time.Millisecond):
			return nil
		}
	})

	if err != nil {
		log.Printf("Error: %v", err)
	}

	// Example 3: Profile a function that returns an error
	fmt.Println("\n=== Example 3: Function with Error Profiling ===")

	err = profiler.ProfileFunction("example_error_function", "main", func() error {
		// Simulate some work
		time.Sleep(25 * time.Millisecond)
		// Return an error to test error tracking
		return fmt.Errorf("example error")
	})

	if err != nil {
		log.Printf("Expected error: %v", err)
	}

	// Example 4: Profile multiple calls to see statistics
	fmt.Println("\n=== Example 4: Multiple Function Calls ===")

	for i := 0; i < 5; i++ {
		err := profiler.ProfileFunction("repeated_function", "main", func() error {
			// Simulate variable work time
			time.Sleep(time.Duration(10+i*5) * time.Millisecond)
			return nil
		})

		if err != nil {
			log.Printf("Error: %v", err)
		}
	}

	// Example 5: Get function profile
	fmt.Println("\n=== Example 5: Get Function Profile ===")

	profile := profiler.GetFunctionProfile("repeated_function")
	if profile != nil {
		fmt.Printf("Function: %s\n", profile.Name)
		fmt.Printf("Call Count: %d\n", profile.CallCount)
		fmt.Printf("Total Duration: %s\n", profiling.FormatDuration(profile.TotalDuration))
		fmt.Printf("Average Duration: %s\n", profiling.FormatDuration(profile.AvgDuration))
		fmt.Printf("Min Duration: %s\n", profiling.FormatDuration(profile.MinDuration))
		fmt.Printf("Max Duration: %s\n", profiling.FormatDuration(profile.MaxDuration))
		fmt.Printf("Total Memory: %s\n", profiling.FormatBytes(profile.TotalMemory))
		fmt.Printf("Max Memory: %s\n", profiling.FormatBytes(profile.MaxMemory))
		fmt.Printf("Error Count: %d\n", profile.ErrorCount)
		fmt.Printf("Last Call: %s\n", profile.LastCallTime.Format(time.RFC3339))
	}

	// Example 6: Get all profiles
	fmt.Println("\n=== Example 6: All Profiles ===")

	allProfiles := profiler.GetAllProfiles()
	fmt.Printf("Total functions profiled: %d\n", len(allProfiles))

	for name, profile := range allProfiles {
		fmt.Printf("\nFunction: %s\n", name)
		fmt.Printf("  Calls: %d, Avg Duration: %s, Errors: %d\n",
			profile.CallCount,
			profiling.FormatDuration(profile.AvgDuration),
			profile.ErrorCount)
	}

	// Example 7: Memory-intensive function
	fmt.Println("\n=== Example 7: Memory-Intensive Function ===")

	err = profiler.ProfileFunction("memory_intensive_function", "main", func() error {
		// Allocate some memory
		data := make([]byte, 1024*1024) // 1MB
		for i := range data {
			data[i] = byte(i % 256)
		}
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	if err != nil {
		log.Printf("Error: %v", err)
	}

	// Example 8: Database-like function simulation
	fmt.Println("\n=== Example 8: Database Function Simulation ===")

	err = profiler.ProfileFunction("database_query", "database", func() error {
		// Simulate database query
		time.Sleep(200 * time.Millisecond)
		// Simulate memory allocation for result
		_ = make([]string, 1000)
		return nil
	})

	if err != nil {
		log.Printf("Error: %v", err)
	}

	// Example 9: API call simulation
	fmt.Println("\n=== Example 9: API Call Simulation ===")

	err = profiler.ProfileFunction("api_call", "external", func() error {
		// Simulate external API call
		time.Sleep(150 * time.Millisecond)
		return nil
	})

	if err != nil {
		log.Printf("Error: %v", err)
	}

	// Example 10: File operation simulation
	fmt.Println("\n=== Example 10: File Operation Simulation ===")

	err = profiler.ProfileFunction("file_operation", "filesystem", func() error {
		// Simulate file I/O
		time.Sleep(75 * time.Millisecond)
		// Simulate memory allocation for file content
		_ = make([]byte, 512*1024) // 512KB
		return nil
	})

	if err != nil {
		log.Printf("Error: %v", err)
	}

	// Final summary
	fmt.Println("\n=== Final Summary ===")

	allProfiles = profiler.GetAllProfiles()
	fmt.Printf("Total functions profiled: %d\n", len(allProfiles))

	var totalCalls int64
	var totalErrors int64
	var totalDuration time.Duration

	for _, profile := range allProfiles {
		totalCalls += profile.CallCount
		totalErrors += profile.ErrorCount
		totalDuration += profile.TotalDuration
	}

	fmt.Printf("Total function calls: %d\n", totalCalls)
	fmt.Printf("Total errors: %d\n", totalErrors)
	fmt.Printf("Total execution time: %s\n", profiling.FormatDuration(totalDuration))
	fmt.Printf("Error rate: %.2f%%\n", float64(totalErrors)/float64(totalCalls)*100)

	fmt.Println("\nProfiling demo completed!")
	fmt.Println("You can now access profiling endpoints at /debug/profiles")
	fmt.Println("And Prometheus metrics at /metrics")
}
