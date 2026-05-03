package main

import (
	"fmt"
)

// -------------------------------------------------------------------------
// ⚠️ CUSTOM AGENT EXCEPTIONS (ERROR HANDLING)
// -------------------------------------------------------------------------

// AgentHallucinationError is a custom exception extending standard errors.
type AgentHallucinationError struct {
	Message string
}

func (e *AgentHallucinationError) Error() string {
	return fmt.Sprintf("AgentHallucinationError: %s", e.Message)
}

// ValidateAgentOutput acts as a validator that throws a custom exception
// if the AI model generates unsafe or unverified data.
func ValidateAgentOutput(output string) error {
	if output == "I don't know but here is a random guess." {
		return &AgentHallucinationError{Message: "Generated response lacks factual backing and violates safety constraints."}
	}
	return nil
}

// DemonstrateCustomErrors shows how custom exceptions are caught and handled.
func DemonstrateCustomErrors(traceID string) {
	fmt.Println("\n==================================================")
	fmt.Println("             CUSTOM AGENT EXCEPTIONS")
	fmt.Println("==================================================")

	LogEvent(traceID, "VALIDATION", "Evaluating LLM output for hallucinations...")
	
	err := ValidateAgentOutput("I don't know but here is a random guess.")
	if err != nil {
		LogEvent(traceID, "EXCEPTION_CAUGHT", fmt.Sprintf("Caught custom exception: %v", err))
		LogEvent(traceID, "EXCEPTION_HANDLE", "Triggering fallback safe response protocol.")
	} else {
		LogEvent(traceID, "EXCEPTION_PASS", "Agent output is valid.")
	}
}
