package main

import (
	"fmt"
	"strings"
)

// -------------------------------------------------------------------------
// 🤖 SINGLE AGENT WORKFLOW
// -------------------------------------------------------------------------

// SingleAgentLoop simulates a basic single agent that Perceives, Reasons, Acts, and Learns.
func SingleAgentLoop(traceID string, environmentInput string) {
	fmt.Println("\n==================================================")
	fmt.Println("             SINGLE AGENT WORKFLOW")
	fmt.Println("==================================================")

	goal := "Find information about CodePathIndia"
	LogEvent(traceID, "GOAL", goal)

	// 1. PERCEIVE
	// Agent gets data from the environment (sensors, user input, APIs, files)
	perception := fmt.Sprintf("Perceived input: %s", environmentInput)
	LogEvent(traceID, "PERCEIVE", perception)

	// 2. REASON
	// Agent analyzes the data, understands context, and decides action
	var action string
	if strings.Contains(strings.ToLower(environmentInput), "codepathindia") {
		action = "Search more about CodePathIndia."
	} else {
		action = "Ask user for CodePathIndia context."
	}
	LogEvent(traceID, "REASON", fmt.Sprintf("Decided action based on context: %s", action))

	// 3. ACT
	// Agent executes the action using tools
	result := "Found relevant information."
	if action != "Search more about CodePathIndia." {
		result = "Waiting for context."
	}
	LogEvent(traceID, "ACT", fmt.Sprintf("Action Result: %s", result))

	// 4. LEARN
	// Agent observes the outcome and learns from feedback
	var learnFeedback string
	if result == "Found relevant information." {
		learnFeedback = "success"
		LogEvent(traceID, "LEARN", fmt.Sprintf("Goal Achieved! Feedback: %s", learnFeedback))
	} else {
		learnFeedback = "continue"
		LogEvent(traceID, "LEARN", fmt.Sprintf("Goal not yet achieved. Feedback: %s", learnFeedback))
	}
}
