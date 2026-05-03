package main

import (
	"fmt"
)

// -------------------------------------------------------------------------
// 🧠 AGENT REFLECTION & SELF-CORRECTION LAYER
// -------------------------------------------------------------------------

// ReflectOnOutcome allows the agent to review its actions and outcomes,
// learn from mistakes, and improve future decisions.
func ReflectOnOutcome(traceID string, action string, outcome string) string {
	LogEvent(traceID, "REFLECT", fmt.Sprintf("Analyzing Outcome -> Action: '%s', Result: '%s'", action, outcome))

	var reflection string
	if outcome == "Failed: Irrelevant Results" {
		reflection = "The search query was not specific enough. Hallucination risk detected."
	} else if outcome == "Success" {
		reflection = "Context was highly relevant. Strategy is optimal."
	} else {
		reflection = "Outcome ambiguous. Requires further verification."
	}

	LogEvent(traceID, "LEARN", fmt.Sprintf("Reflection insight: %s", reflection))
	return reflection
}

// ImproveStrategy applies the new strategy based on the reflection
func ImproveStrategy(traceID string, reflection string) string {
	improvedAction := "Use more specific keywords and apply metadata filters."
	LogEvent(traceID, "IMPROVE", fmt.Sprintf("Next Action Adjusted: %s", improvedAction))
	return improvedAction
}

// SelfCorrectingAgent Loop
// Continually evaluates actions, learns from mistakes, and improves future decisions.
func SelfCorrectingExecution(traceID string, goal string) bool {
	LogEvent(traceID, "INFO", fmt.Sprintf("Self-Correcting Loop Started for Goal: %s", goal))
	
	goalAchieved := false
	attempt := 1

	for !goalAchieved && attempt <= 3 {
		fmt.Printf("\n--- [Attempt %d] ---\n", attempt)
		
		// 1. Plan & Act
		action := "Search Vector DB for standard query"
		LogEvent(traceID, "ACT", action)

		// Simulated failure on attempt 1
		success := false
		outcome := "Failed: Irrelevant Results"
		if attempt > 1 {
			success = true
			outcome = "Success"
		}

		// 2. Observe & Evaluate
		LogEvent(traceID, "OBSERVE", outcome)

		if success {
			goalAchieved = true
			LogEvent(traceID, "SUCCESS", "Goal Achieved!")
		} else {
			// 3. Detect & Correct
			LogEvent(traceID, "DETECT", "Plan failed. Initiating Reflection & Correction...")
			
			// Reflection Phase
			reflection := ReflectOnOutcome(traceID, action, outcome)
			
			// 4. Learn & Improve
			newPlan := ImproveStrategy(traceID, reflection)
			LogEvent(traceID, "CORRECT", fmt.Sprintf("Adjusting plan to: %s", newPlan))
		}
		attempt++
	}

	return goalAchieved
}
