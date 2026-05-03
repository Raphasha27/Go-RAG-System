package main

import (
	"crypto/rand"
	"fmt"
	"time"
)

// -------------------------------------------------------------------------
// 👁️ AGENT OBSERVABILITY LAYER
// -------------------------------------------------------------------------

// generateTraceID creates a simple pseudo-UUID for tracking request flows
func generateTraceID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// LogEvent captures important events, inputs, outputs, and errors for auditing
func LogEvent(traceID, level, message string) {
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Printf("[%s] [TraceID: %s] [%s] %s\n", timestamp, traceID, level, message)
}

// ObserveExecution wraps a function execution with Trace, Log, and Metric pillars
func ObserveExecution(traceID string, stepName string, fn func() error) error {
	LogEvent(traceID, "INFO", fmt.Sprintf("Starting step: %s", stepName))
	
	start := time.Now()
	err := fn()
	duration := time.Since(start)

	if err != nil {
		LogEvent(traceID, "ERROR", fmt.Sprintf("Step '%s' failed after %v. Error: %v", stepName, duration, err))
		return err
	}

	LogEvent(traceID, "SUCCESS", fmt.Sprintf("Step '%s' completed successfully in %v", stepName, duration))
	return nil
}
