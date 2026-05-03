package main

import "fmt"

// -------------------------------------------------------------------------
// 🛡️ AGENT SECURITY LAYER
// -------------------------------------------------------------------------

// isAuthorized checks whether an AI agent (or user) is authorized to perform a specific action based on their role.
// This implements the Privilege Limitation and Authorization best practices.
func isAuthorized(role string, action string) bool {
	// Define Role-Based Access Control (RBAC) for the Agentic System
	permissions := map[string][]string{
		"admin": {"read", "write", "delete", "execute"},
		"user":  {"read", "write"},
		"agent": {"read", "execute"},
	}

	allowedActions, exists := permissions[role]
	if !exists {
		return false // Role does not exist, default deny
	}

	for _, allowedAction := range allowedActions {
		if allowedAction == action {
			return true
		}
	}
	return false
}

// EnforceSecurity validates tool execution attempts by the Agent
func EnforceSecurity(role string, toolName string) bool {
	// Map specific tools to required authorization actions
	requiredAction := "execute"
	if toolName == "delete_record" {
		requiredAction = "delete"
	} else if toolName == "write_file" {
		requiredAction = "write"
	}

	// Audit Logging & Authorization Check
	if isAuthorized(role, requiredAction) {
		fmt.Printf("[🔒 Security Audit]: Role '%s' AUTHORIZED to perform '%s' via tool '%s'.\n", role, requiredAction, toolName)
		return true
	}

	fmt.Printf("[❌ Security Alert]: Role '%s' DENIED execution of tool '%s'.\n", role, toolName)
	return false
}
