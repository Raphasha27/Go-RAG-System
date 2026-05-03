package main

import (
	"fmt"
	"sync"
	"time"
)

// -------------------------------------------------------------------------
// ⚡ CONCURRENCY IN AI AGENTS
// -------------------------------------------------------------------------

// DemonstrateConcurrency shows how an Agentic system handles multiple tasks 
// simultaneously using Goroutines (Go's lightweight threads).
func DemonstrateConcurrency(traceID string) {
	fmt.Println("\n==================================================")
	fmt.Println("             AGENT CONCURRENCY (GOROUTINES)")
	fmt.Println("==================================================")

	LogEvent(traceID, "CONCURRENCY_START", "Launching multiple agent sub-tasks concurrently to optimize performance...")

	var wg sync.WaitGroup
	tasks := []string{"Fetch Real-time Weather API", "Query Vector Database", "Update Long-Term Memory"}

	for _, task := range tasks {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			LogEvent(traceID, "CONCURRENCY_RUNNING", fmt.Sprintf("Executing asynchronously: '%s'", t))
			time.Sleep(300 * time.Millisecond) // Simulate I/O bound work
			LogEvent(traceID, "CONCURRENCY_DONE", fmt.Sprintf("Completed: '%s'", t))
		}(task)
	}

	wg.Wait()
	LogEvent(traceID, "CONCURRENCY_SUCCESS", "All concurrent threads synchronized and completed. Race conditions avoided.")
}
