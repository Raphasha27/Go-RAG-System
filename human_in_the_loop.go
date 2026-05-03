package main

import (
	"fmt"
	"strings"
)

// -------------------------------------------------------------------------
// 🙋 HUMAN-IN-THE-LOOP (HITL) LAYER
// -------------------------------------------------------------------------

// HumanApproval requests human feedback before an AI agent executes a critical or destructive action.
// It combines AI autonomy with human judgment to ensure accuracy, reliability, and safety.
func HumanApproval(traceID string, proposedAction string) string {
	fmt.Println("\n==================================================")
	fmt.Println("             HUMAN-IN-THE-LOOP REVIEW")
	fmt.Println("==================================================")
	LogEvent(traceID, "HITL_REQUEST", fmt.Sprintf("Critical action requires human review: '%s'", proposedAction))

	fmt.Printf("\n[AI Agent Proposes]: %s\n", proposedAction)
	fmt.Print("Review the above response. Approve (A) / Modify (M) / Reject (R): ")

	// Note: In a real interactive application, you would use bufio.NewReader(os.Stdin).
	// For this fully automated architectural blueprint, we simulate the human input.
	humanInput := "A" // Simulated human response
	fmt.Printf("%s\n", humanInput)

	finalAction := proposedAction

	switch strings.ToUpper(humanInput) {
	case "A":
		LogEvent(traceID, "HITL_APPROVE", "Human approved the action.")
		finalAction = proposedAction + " (Approved by Human)"
	case "M":
		modifiedInput := "Extract vector context but limit to the top 1 result."
		fmt.Printf("Enter your modified response: %s\n", modifiedInput)
		LogEvent(traceID, "HITL_MODIFY", fmt.Sprintf("Human modified action to: '%s'", modifiedInput))
		finalAction = modifiedInput
	case "R":
		reason := "This query is too broad. Please narrow down."
		fmt.Printf("Enter the reason for rejection: %s\n", reason)
		LogEvent(traceID, "HITL_REJECT", fmt.Sprintf("Human rejected action. Reason: '%s'", reason))
		finalAction = "Task Rejected. Reason: " + reason
	default:
		LogEvent(traceID, "HITL_ERROR", "Invalid input.")
		finalAction = "Invalid input. No action taken."
	}

	LogEvent(traceID, "ACT", fmt.Sprintf("Final Action taken: %s", finalAction))
	return finalAction
}
